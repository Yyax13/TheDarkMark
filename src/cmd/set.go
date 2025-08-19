package cmd

import (
	"fmt"
	"strings"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

var set types.Command = types.Command{
	Name: "set",
	Description: "Sets an option",
	Listable: false,
	HelpListDescription: "",
	Run: setCommand,

}

func setCommand(mainEnv *types.MainEnvType, args []string) {
	if len(args) < 3 {
		misc.PanicWarn("Usage: set <option> <value>\n", true)
		return

	}

	rawOptionsName := strings.TrimSpace(args[1])
	newVal := strings.TrimSpace(args[2])
	_, exists := mainEnv.Module.Options[rawOptionsName]
	if !exists {
		misc.PanicWarn(fmt.Sprintf("Option %s not found, use options to view avaliable options in current module\n", args[1]), true)
		return

	}

	mainEnv.Module.SetOptionVal(rawOptionsName, newVal)

}

func init() {
	RegisterNewCommand(&set)

}
