package types

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/misc"
)

type Chamber struct {
	Name        string
	Description string
	Runes     map[string]*Rune
	Parallel    bool
	Execute     func(runes map[string]*Rune)
}

func (chb *Chamber) ListAvaliableRunes() {
	if len(chb.Runes) == 0 {
		misc.PanicWarn(fmt.Sprintf("The command %v haven't runes", chb.Name), false) // needs to be "incantation" and not "command"
		return

	}

	fmt.Printf("Avaliable runes for %v:\n", chb.Name)
	fmt.Printf("	%-10s %-15s %-10s %-15s\n", "Name", "Current Value", "Required", "Description")
	for _, opt := range chb.Runes {
		req := "No"
		if opt.Required {
			req = "Yes"

		}

		fmt.Printf("	%-10s %-15s %-10s %-15s\n", opt.Name, fmt.Sprintf("%v", opt.Value), req, opt.Description)

	}

	fmt.Print("\n")

}

func (chb *Chamber) SetRuneVal(runeName string, RuneVal string) error {
	opt, exists := chb.Runes[runeName]
	if !exists {
		return fmt.Errorf("not found the option %s", opt.Name) // todo: change option to rune

	}

	opt.Value = RuneVal
	misc.SysLog(fmt.Sprintf("Successfuly set %s value to %v\n\n", opt.Name, opt.Value), false)
	return nil

}

func (chb *Chamber) CheckRunesValue() error {
	for _, opt := range chb.Runes {
		if opt.Required && opt.Value == "" {
			return fmt.Errorf("rune %v is required and not set by user", opt.Name)

		}

	}

	return nil

}
