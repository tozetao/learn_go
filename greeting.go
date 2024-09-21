package main

import (
	"errors"
)

func Hello(name string) (string, error) {
	if name == "" {
		return name, errors.New("empty name")
	}
	message := "hi, " + name
	return message, nil
}
