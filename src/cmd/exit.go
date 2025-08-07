package cmd

import (
	"os"
	"github.com/Yyax13/onTop-C2/src/types"

)

var exit types.Command = types.Command{
	Name: "exit",
	Description: "Exit the program without any messages",
	Listable: false,
	HelpListDescription: "",
	Run: exitCommand,
	
}

func exitCommand(mainEnv *types.MainEnvType, _ []string) {
	os.Exit(0)

}

func init() {
	RegisterNewCommand(&exit)

}