package modules

import "gopkg.in/telebot.v3"

type Command struct {
	// Name of the command
	Name interface{}

	// Func that takes in a telebot.Context as the only arg and returns an error
	Handler func(telebot.Context) error
}
type Listener struct {
	// TODO:
}

type Data struct {
	// Returns a slice of commands
	Commands *[]Command

	// Returns a slice of listeners
	Listeners *[]Listener
}
type Module interface {
	// A func that returns a slice of commands
	Init() *Data
}
