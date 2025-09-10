package misc

import (
	"errors"
	"fmt"
)

var AvaliableColors map[string]string = map[string]string{
	"red":          "\x1b[38;2;215;0;0m",
	"green":        "\x1b[38;2;0;215;0m",
	"blue":         "\x1b[38;2;0;0;215m",
	"yellow":       "\x1b[38;2;215;215;0m",
	"magenta":      "\x1b[38;2;215;0;215m",
	"cyan":         "\x1b[38;2;0;215;215m",
	"white":        "\x1b[38;2;215;215;215m",
	"black":        "\x1b[38;2;0;0;0m",

	"red_bold":     "\x1b[1;38;2;215;0;0m",
	"green_bold":   "\x1b[1;38;2;0;215;0m",
	"blue_bold":    "\x1b[1;38;2;0;0;215m",
	"yellow_bold":  "\x1b[1;38;2;215;215;0m",
	"magenta_bold": "\x1b[1;38;2;215;0;215m",
	"cyan_bold":    "\x1b[1;38;2;0;215;215m",
	"white_bold":   "\x1b[1;38;2;215;215;215m",
	"black_bold":   "\x1b[1;38;2;0;0;0m",

	"light_red":    "\x1b[38;2;240;0;0m",
	"light_green":  "\x1b[38;2;0;240;0m",
	"light_blue":   "\x1b[38;2;0;0;240m",
	"light_yellow": "\x1b[38;2;240;240;0m",
	"light_magenta":"\x1b[38;2;240;0;240m",
	"light_cyan":   "\x1b[38;2;0;240;240m",
	"light_white":  "\x1b[38;2;240;240;240m",
	"light_black":  "\x1b[38;2;60;60;60m",

	"light_red_bold":    "\x1b[1;38;2;240;0;0m",
	"light_green_bold":  "\x1b[1;38;2;0;240;0m",
	"light_blue_bold":   "\x1b[1;38;2;0;0;240m",
	"light_yellow_bold": "\x1b[1;38;2;240;240;0m",
	"light_magenta_bold":"\x1b[1;38;2;240;0;240m",
	"light_cyan_bold":   "\x1b[1;38;2;0;240;240m",
	"light_white_bold":  "\x1b[1;38;2;240;240;240m",
	"light_black_bold":  "\x1b[1;38;2;60;60;60m",

	"dark_red":     "\x1b[38;2;110;0;0m",
	"dark_green":   "\x1b[38;2;0;110;0m",
	"dark_blue":    "\x1b[38;2;0;0;110m",
	"dark_yellow":  "\x1b[38;2;110;110;0m",
	"dark_magenta": "\x1b[38;2;110;0;110m",
	"dark_cyan":    "\x1b[38;2;0;110;110m",
	"dark_white":   "\x1b[38;2;110;110;110m",
	"dark_black":   "\x1b[38;2;0;0;0m",

	"dark_red_bold":     "\x1b[1;38;2;110;0;0m",
	"dark_green_bold":   "\x1b[1;38;2;0;110;0m",
	"dark_blue_bold":    "\x1b[1;38;2;0;0;110m",
	"dark_yellow_bold":  "\x1b[1;38;2;110;110;0m",
	"dark_magenta_bold": "\x1b[1;38;2;110;0;110m",
	"dark_cyan_bold":    "\x1b[1;38;2;0;110;110m",
	"dark_white_bold":   "\x1b[1;38;2;110;110;110m",
	"dark_black_bold":   "\x1b[1;38;2;0;0;0m",

	"reset":        "\x1b[0m",
}

func Colors(text, color string) (string, error) {
	if _, exists := AvaliableColors[color]; !exists {
		return "", errors.New("color not avaliable") // idk why this handle exists, just the devs ll use that and i think that they know the avaliable colors

	}

	return fmt.Sprintf("%s%v%s", AvaliableColors[color], text, AvaliableColors["reset"]), nil

}

