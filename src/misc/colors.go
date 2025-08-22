package misc

import (
	"errors"
	"fmt"
)

var avaliableColors map[string]string = map[string]string{
    "red":          "\033[38;2;255;0;0m",
    "green":        "\033[38;2;0;255;0m",
    "blue":         "\033[38;2;0;0;255m",
    "yellow":       "\033[38;2;255;255;0m",
    "magenta":      "\033[38;2;255;0;255m",
    "cyan":         "\033[38;2;0;255;255m",
    "white":        "\033[38;2;255;255;255m",
    "black":        "\033[38;2;0;0;0m",

    "red_bold":     "\033[1;38;2;255;0;0m",
    "green_bold":   "\033[1;38;2;0;255;0m",
    "blue_bold":    "\033[1;38;2;0;0;255m",
    "yellow_bold":  "\033[1;38;2;255;255;0m",
    "magenta_bold": "\033[1;38;2;255;0;255m",
    "cyan_bold":    "\033[1;38;2;0;255;255m",
    "white_bold":   "\033[1;38;2;255;255;255m",
    "black_bold":   "\033[1;38;2;0;0;0m",

    "light_red":    "\033[38;2;255;70;70m",
    "light_green":  "\033[38;2;70;255;70m",
    "light_blue":   "\033[38;2;70;70;255m",
    "light_yellow": "\033[38;2;255;255;70m",
    "light_magenta":"\033[38;2;255;70;255m",
    "light_cyan":   "\033[38;2;70;255;255m",
    "light_white":  "\033[38;2;255;255;255m",
    "light_black":  "\033[38;2;70;70;70m",

    "light_red_bold":    "\033[1;38;2;255;70;70m",
    "light_green_bold":  "\033[1;38;2;70;255;70m",
    "light_blue_bold":   "\033[1;38;2;70;70;255m",
    "light_yellow_bold": "\033[1;38;2;255;255;70m",
    "light_magenta_bold":"\033[1;38;2;255;70;255m",
    "light_cyan_bold":   "\033[1;38;2;70;255;255m",
    "light_white_bold":  "\033[1;38;2;255;255;255m",
    "light_black_bold":  "\033[1;38;2;70;70;70m",

    "dark_red":     "\033[38;2;185;0;0m",
    "dark_green":   "\033[38;2;0;185;0m",
    "dark_blue":    "\033[38;2;0;0;185m",
    "dark_yellow":  "\033[38;2;185;185;0m",
    "dark_magenta": "\033[38;2;185;0;185m",
    "dark_cyan":    "\033[38;2;0;185;185m",
    "dark_white":   "\033[38;2;185;185;185m",
    "dark_black":   "\033[38;2;0;0;0m",

    "dark_red_bold":     "\033[1;38;2;185;0;0m",
    "dark_green_bold":   "\033[1;38;2;0;185;0m",
    "dark_blue_bold":    "\033[1;38;2;0;0;185m",
    "dark_yellow_bold":  "\033[1;38;2;185;185;0m",
    "dark_magenta_bold": "\033[1;38;2;185;0;185m",
    "dark_cyan_bold":    "\033[1;38;2;0;185;185m",
    "dark_white_bold":   "\033[1;38;2;185;185;185m",
    "dark_black_bold":   "\033[1;38;2;0;0;0m",

    "reset":        "\033[0m",
}

func Colors(text, color string) (string, error) {
	if _, exists := avaliableColors[color]; !exists {
		return "", errors.New("color not avaliable")

	}

	return fmt.Sprintf("%s%v%s", avaliableColors[color], text, avaliableColors["reset"]), nil

}

