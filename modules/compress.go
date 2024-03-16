package modules

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
	"gopkg.in/telebot.v3"
)

type Compress struct {
	Module
}

func downloadVideo(url, destination string) error {
	// Make a GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the output file
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomString(length int) string {
	// Make a string of all the characters
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	// Make a byte array of the characters
	charactersLength := len(characters)
	charactersBytes := []byte(characters)

	// Make a byte array of the random characters
	randomBytes := make([]byte, length)

	// Generate random bytes
	for i := range randomBytes {
		randomBytes[i] = charactersBytes[rand.Intn(charactersLength)]
	}

	// Return the random string
	return string(randomBytes)
}

func CompressVideo(filepath string) error {
	// Compress the video
	cmd := exec.Command("ffmpeg", "-i", filepath, "-vcodec", "h264", "-acodec", "mp3", "-crf", "28", filepath+".mp4")
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Delete the old video
	err = os.Remove(filepath)
	if err != nil {
		return err
	}

	// Rename the new video
	err = os.Rename(filepath+".mp4", filepath)
	if err != nil {
		return err
	}

	return nil
}

func CompressHandler(ctx telebot.Context) error {
	// Get the video
	video := ctx.Message().Video
	log.Info(video.FileURL)
	// Download the video
	err := downloadVideo(video.FileURL, video.FileName)
	if err != nil {
		ctx.Reply("Error: " + err.Error())
		return nil
	}

	// Compress the video
	err = CompressVideo(video.FileName)
	if err != nil {
		ctx.Reply("Error: " + err.Error())
		return nil
	}

	// Reply with the video
	// Get the precentage of the compression from the original file size to the new file size
	originalFileSize := video.FileSize
	newFileSize, err := os.Stat(video.FileName)
	if err != nil {
		ctx.Reply("Error: " + err.Error())
		return nil
	}
	percentage := (float64(newFileSize.Size()) / float64(originalFileSize)) * 100

	ctx.Reply(&telebot.Video{File: telebot.FromDisk(video.FileName), Caption: fmt.Sprintf("Compressed by %.2f%%\n\nOriginal: %d\n\nNew: %d", percentage, originalFileSize, newFileSize.Size()), FileName: generateRandomString(10) + ".mp4"})

	// Delete the video
	err = os.Remove(video.FileName)
	if err != nil {
		ctx.Reply("Error: " + err.Error())
	}

	return nil
}

func (m *Compress) Init() *Data {
	return &Data{
		Commands: &[]Command{
			{
				Name:    telebot.OnVideo,
				Handler: CompressHandler,
			},
		},
	}
}
