package modules

import (
	"gopkg.in/telebot.v3"
	"strings"
)

type Carny struct {
	Module
}

func translateToCarny(englishWord string) string {
	vowels := []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U"}
	syllables := strings.Split(englishWord, "")

	translatedWord := ""

	for _, syllable := range syllables {
		if len(syllable) == 0 {
			continue
		}

		startsWithVowel := false
		for _, vowel := range vowels {
			if strings.HasPrefix(syllable, vowel) {
				startsWithVowel = true
				break
			}
		}

		if startsWithVowel {
			translatedWord += "earz" + syllable
		} else {
			translatedWord += syllable
		}
	}

	return translatedWord
}

func translateToEnglish(spokenWord string) string {
	englishWord := ""
	syllables := strings.Split(spokenWord, "earz")

	for _, syllable := range syllables {
		if len(syllable) == 0 {
			continue
		}

		englishWord += syllable
	}

	return englishWord
}

func CarnyHandler(ctx telebot.Context) error {
	// Concat all the args seperated by spaces
	args := ctx.Args()
	var message string
	for _, arg := range args {
		message += arg + " "
	}

	var translatedMessage string
	if strings.Contains(message, "earz") {
		translatedMessage = translateToEnglish(message)
	} else {
		translatedMessage = translateToCarny(message)
	}
	ctx.Reply(translatedMessage)

	return nil
}

func (m *Carny) Init(_ *telebot.Bot) *Data {
	// Return a slice of commands
	return &Data{
		Commands: &[]Command{
			{
				Name:    "/carny",
				Handler: CarnyHandler,
			},
		},
	}
}
