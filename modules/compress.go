package modules

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/dustin/go-humanize"
	"gopkg.in/telebot.v3"
)

type Compress struct {
	Module
}

var (
	allowedDocExtensions     = []string{".mp4", ".mov", ".avi", ".mkv", ".flv", ".wmv", ".webm"}
	alwaysCompressExtensions = []string{".mkv", ".webm"} // These will always be compressed regardless of size due to them not embedding on telegram
)

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
	cmd := exec.Command("ffmpeg", "-i", filepath, "-vcodec", "h264", "-acodec", "mp3", "-crf", "28", filepath+".out.mp4") // Perhaps use different codecs?
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func CompressHandler(b *telebot.Bot) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		// This is so fucked but I am not sure how else to do it...
		var isVideo bool
		var video *telebot.Video
		var document *telebot.Document
		var forceCompress bool = false
		video = ctx.Message().Video
		if video != nil {
			isVideo = true
		} else {
			document = ctx.Message().Document
			if document != nil {
				// Check if the file name ends with an allowed extension
				allowed := false
				for _, ext := range allowedDocExtensions {
					if ext == filepath.Ext(document.FileName) {
						allowed = true
						break
					}
				}
				for _, ext := range alwaysCompressExtensions {
					if ext == filepath.Ext(document.FileName) {
						forceCompress = true
						break
					}
				}
				if !allowed {
					return nil
				}
			}
		}
		var fileName string
		var rawFileSize int64

		// Uhhhh
		if isVideo {
			fileName = video.FileName
			rawFileSize = video.FileSize
		} else {
			fileName = document.FileName
			rawFileSize = document.FileSize
		}
		// Check if the file is larger than 8MB
		if rawFileSize < 8000000 && !forceCompress {
			log.Infof("File: %s is %s, skipping compression", fileName, humanize.Bytes(uint64(rawFileSize)))
			return nil // We don't want to compress files that are already small
		}
		// React to the message
		err := ctx.Notify(telebot.Typing)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
		log.Infof("File: %s: Downloading... | File Size: %v", fileName, rawFileSize)

		// Download the file
		tempName := generateRandomString(10)
		if isVideo {
			err = b.Download(&video.File, "/tmp/"+tempName)
		} else {
			err = b.Download(&document.File, "/tmp/"+tempName)
		}
		if err != nil {
			err = ctx.Reply("Error: " + err.Error())
			if err != nil {
				log.Error(err.Error())
			}
			return nil
		}

		// Compress the file
		err = CompressVideo("/tmp/" + tempName)
		if err != nil {
			err = ctx.Reply("Error: " + err.Error())
			if err != nil {
				log.Error(err.Error())
			}
			return nil
		}

		// Reply with the compressed file
		compressedFileName := tempName + ".out.mp4"
		newFile, err := os.Stat("/tmp/" + compressedFileName)
		if err != nil {
			err = ctx.Reply("Error: " + err.Error())
			if err != nil {
				log.Error(err.Error())
			}
			return nil
		}

		// Calculate compression percentage
		percentage := (float64(rawFileSize-newFile.Size()) / float64(rawFileSize)) * 100

		log.Infof("File: %s: Compressed by %.2f%% | New File Size: %v", fileName, percentage, newFile.Size())

		err = ctx.Reply(&telebot.Video{File: telebot.FromDisk("/tmp/" + compressedFileName), Caption: fmt.Sprintf("Compressed by %.2f%%\n\nOriginal: %s\n\nNew: %s", percentage, humanize.Bytes(uint64(rawFileSize)), humanize.Bytes(uint64(newFile.Size()))), FileName: tempName + ".mp4"})
		if err != nil {
			log.Error(err.Error())
		}

		// Delete the compressed file
		err = os.Remove("/tmp/" + compressedFileName)
		if err != nil {
			err = ctx.Reply("Error: " + err.Error())
			if err != nil {
				log.Error(err.Error())
			}
		}

		// Delete the original file
		err = os.Remove("/tmp/" + tempName)
		if err != nil {
			err = ctx.Reply("Error: " + err.Error())
			if err != nil {
				log.Error(err.Error())
			}
		}

		return nil
	}
}

func (m *Compress) Init(b *telebot.Bot) *Data {
	return &Data{
		Commands: &[]Command{
			{
				Name:    telebot.OnVideo,
				Handler: CompressHandler(b), // Better way to do this??
			},
			{
				Name:    telebot.OnDocument,
				Handler: CompressHandler(b),
			},
		},
	}
}
