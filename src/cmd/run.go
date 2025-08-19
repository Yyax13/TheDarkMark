package cmd

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/types"
	"github.com/Yyax13/onTop-C2/src/misc"

)

var run types.Command = types.Command{
	Name: "run",
	Description: "Runs the current module",
	Listable: false,
	HelpListDescription: "",
	Run: runCommand,

}

func runCommand(mainEnv *types.MainEnvType, args []string) {
	err := mainEnv.Module.CheckOptionsValue()
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("%v\n", err.Error()), true)
		return

	}


	switch mainEnv.Module.Parallel {
	case true:
		go mainEnv.Module.Execute(mainEnv.Module.Options)
		return

	default:
		mainEnv.Module.Execute(mainEnv.Module.Options)
		return

	}

}

func init() {
	RegisterNewCommand(&run)

}
