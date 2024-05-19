package cmd

import (
	"Divine-Intellect/modules"
	"os"
	"time"

	"github.com/charmbracelet/log"
	tele "gopkg.in/telebot.v3"
)

// Create an array of Modules
var registeredModules = []modules.Module{
	&modules.Hello{},
	&modules.Downloader{},
	&modules.Compress{},
	&modules.Carny{},
}

func loadModules(b *tele.Bot) {
	for _, module := range registeredModules {
		// Get the commands from the module
		data := module.Init(b)
		if data == nil {
			continue
		}

		if data.Commands == nil {
			continue
		}

		// Register the commands
		for _, command := range *data.Commands {
			b.Handle(command.Name, command.Handler)
		}
	}
}

func Start() {
	log.Info("Starting bot...")
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info("Loading modules...")
	loadModules(b)
	log.Info("Starting bot...")

	b.Start()
}
