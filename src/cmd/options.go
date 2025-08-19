package cmd

import (
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

var options types.Command = types.Command{
	Name:                "options",
	Description:         "Show the options in current module",
	Listable:            false,
	HelpListDescription: "",
	Run:                 optionsRun,
}

func optionsRun(mainEnv *types.MainEnvType, args []string) {
	if mainEnv.Module.Name == "" {
		misc.PanicWarn("You aren't in a module, run list use to view avaliable modules.\n", true)
		return

	}

	mainEnv.Module.ListAvaliableOptions()

}

func init() {
	RegisterNewCommand(&options)

}
