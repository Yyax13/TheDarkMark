package cmd

import (
	"fmt"
	"os"
	"github.com/Yyax13/onTop-C2/src/types"
	"github.com/Yyax13/onTop-C2/src/misc"

)

var clear types.Command = types.Command{
	Name: "clear",
	Description: "Clear the stdout and print the banner",
	Listable: false,
	HelpListDescription: "",
	Run: clearCommand,

}

func clearCommand(mainEnv *types.MainEnvType, _ []string) {
	e := misc.ForceClearStdout()
	if e != nil {
		misc.PanicWarn(fmt.Sprintf("\nAn error ocurred during stdout forced clear: %v", e), true)
		os.Exit(0)

	}

	misc.PrintBanner()

}

func init() {
	RegisterNewCommand(&clear)

}