package modules

import "gopkg.in/telebot.v3"

// Inherit Module

type Hello struct {
	Module
}

func HelloHandler(ctx telebot.Context) error {
	// Send a message to the user
	ctx.Reply("Hello, world!")
	return nil
}

func (m *Hello) Init() *Data {
	// Return a slice of commands
	return &Data{
		Commands: &[]Command{
			{
				Name:    "/hello",
				Handler: HelloHandler,
			},
		},
	}
}
