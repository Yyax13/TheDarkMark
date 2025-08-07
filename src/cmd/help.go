package cmd

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/types"

)

var help types.Command = types.Command{
	Name: "help",
	Description: "Shows all avaliable commands",
	Listable: false,
	HelpListDescription: "",
	Run: helpCommand,

}

func helpCommand(mainEnv *types.MainEnvType, _ []string) {
	fmt.Println("Avaliable commands:")
	for _, cmd := range AvaliableCommands {
		fmt.Printf("	%-13s: %-16s", cmd.Name, cmd.Description)

	}

}

func init() {
	RegisterNewCommand(&help)

}