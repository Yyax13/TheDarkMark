package misc

import (
	"fmt"
	"errors"

)

var avaliableColors map[string]string = map[string]string{
	"black":   "\033[1;30m",
	"red":     "\033[1;31m",
	"green":   "\033[1;32m",
	"yellow":  "\033[1;33m",
	"blue":    "\033[1;34m",
	"magenta": "\033[1;35m",
	"cyan":    "\033[1;36m",
	"white":   "\033[1;37m",
	"reset":   "\033[0m",

}

func Colors(text, color string) (string, error) {
	if _, exists := avaliableColors[color]; !exists {
		return "", errors.New("color not avaliable")

	}

	return fmt.Sprintf("%s%v%s", avaliableColors[color], text, avaliableColors["reset"]), nil

}