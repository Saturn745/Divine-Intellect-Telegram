package modules

import (
	"context"

	"gopkg.in/telebot.v3"

	"github.com/wader/goutubedl"
)

type Downloader struct {
	Module
}

// Make an array of strings
var urls = []string{
	"https://cdn.discordapp.com",
	"https://media.discordapp.net",
	"https://images-ext-1.discordapp.net",
	"https://tiktok.com",
	"https://www.tiktok.com",
}

func download(ctx telebot.Context, url string) {
	result, err := goutubedl.New(context.Background(), url, goutubedl.Options{})
	ctx.Reply("Downloading: " + result.Info.Title)
	if err != nil {
		ctx.Reply("Error: " + err.Error())
		return
	}
	downloadResult, err := result.Download(context.Background(), "")
	if err != nil {
		ctx.Reply("Error: " + err.Error())
		return
	}
	ctx.Reply("Uploading: " + result.Info.Title)
	file := &telebot.Video{File: telebot.FromReader(downloadResult), FileName: result.Info.Title + ".mp4"}
	ctx.Reply(file)
}

func DownloaderHandler(ctx telebot.Context) error {
	args := ctx.Args()
	if len(args) == 0 {
		ctx.Reply("Usage: /download <url>")
		return nil
	}
	url := args[0]
	download(ctx, url)
	return nil
}

func AutoDownloadHandler(ctx telebot.Context) error {
	for _, url := range urls {
		// If the message starts with the url
		if len(ctx.Text()) >= len(url) && ctx.Text()[:len(url)] == url {
			download(ctx, ctx.Text())
			return nil
		}
	}
	return nil
}

func (m *Downloader) Init(_ *telebot.Bot) *Data {
	return &Data{
		Commands: &[]Command{
			{
				Name:    "/download",
				Handler: DownloaderHandler,
			},
			{
				Name:    telebot.OnText,
				Handler: AutoDownloadHandler,
			},
		},
	}
}
