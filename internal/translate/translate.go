package translate

import (
	"errors"
	"fmt"
	"strings"
)

var supportedCommands = map[string]any{"export": struct{}{}, "reset": struct{}{}}

func FindCommand(text string) (*Command, error) {
	if text != "" && strings.HasPrefix(text, "/") {
		text = strings.TrimLeft(text, "/")
		tokens := strings.Split(text, " ")
		if len(tokens) > 0 {
			name := strings.TrimSpace(tokens[0])
			if name == "" {
				return nil, errors.New("command is empty")
			}
			if _, ok := supportedCommands[name]; ok {
				args := tokens[1:]
				if (name == "export" || name == "reset") && len(args) == 0 {
					return nil, fmt.Errorf("%s should have lang argument", name)
				}
				return &Command{
					Name: name,
					Args: args,
				}, nil
			}
		}
	}
	return nil, nil
}

// FindPhrase finds phrase to save in the database
func FindPhrase(text string) (*Phrase, error) {
	if text != "" && strings.HasPrefix(text, "/") {
		text = strings.TrimLeft(text, "/")
		tokens := strings.Split(text, "\n")
		if len(tokens) != 3 {
			return nil, errors.New("tokens size is incorrect")
		}
		lang := strings.TrimSpace(tokens[0])
		if lang == "" {
			return nil, errors.New("lang is empty")
		}
		phrase := strings.TrimSpace(tokens[1])
		if phrase == "" {
			return nil, errors.New("phrase is empty")
		}
		translation := strings.TrimSpace(tokens[2])
		if translation == "" {
			return nil, errors.New("translation is empty")
		}
		return &Phrase{
			Lang:        lang,
			Phrase:      phrase,
			Translation: tokens[2],
		}, nil
	}
	return nil, nil
}

// Phrase is the phrase to store in the database
type Phrase struct {
	Lang        string
	Phrase      string
	Translation string
}

type Command struct {
	Name string
	Args []string
}
