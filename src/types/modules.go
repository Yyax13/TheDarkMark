package types

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/misc"
)

type Module struct {
	Name        string
	Description string
	Options     map[string]*Option
	Parallel    bool
	Execute     func(opts map[string]*Option)
}

func (cmd *Module) ListAvaliableOptions() {
	if len(cmd.Options) == 0 {
		misc.PanicWarn(fmt.Sprintf("The command %v haven't options", cmd.Name), false)
		return

	}

	fmt.Printf("Avaliable options for %v:\n", cmd.Name)
	fmt.Printf("	%-10s %-15s %-10s %-15s\n", "Name", "Current Value", "Required", "Description")
	for _, opt := range cmd.Options {
		req := "No"
		if opt.Required {
			req = "Yes"

		}

		fmt.Printf("	%-10s %-15s %-10s %-15s\n", opt.Name, fmt.Sprintf("%v", opt.Value), req, opt.Description)

	}

	fmt.Print("\n")

}

func (cmd *Module) SetOptionVal(optionName string, OptionVal any) error {
	opt, exists := cmd.Options[optionName]
	if !exists {
		return fmt.Errorf("not found the option %s", opt.Name)

	}

	opt.Value = OptionVal
	misc.SysLog(fmt.Sprintf("Successfuly set %s value to %v\n\n", opt.Name, opt.Value), false)
	return nil

}

func (cmd *Module) CheckOptionsValue() error {
	for _, opt := range cmd.Options {
		if opt.Required && opt.Value == nil {
			return fmt.Errorf("option %v is required and not set by user", opt.Name)

		}

	}

	return nil

}
