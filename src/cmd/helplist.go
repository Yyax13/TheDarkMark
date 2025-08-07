package cmd

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/types"
	
)

var helplist types.Command = types.Command{
	Name: "helplist",
	Description: "Shows all commands compatible with list command",
	Listable: false,
	HelpListDescription: "",
	Run: helplistCommands,

}

func helplistCommands(mainEnv *types.MainEnvType, _ []string) {
	fmt.Println("List able commands:")
	for _, cmd := range AvaliableCommands {
		if cmd.Listable && cmd.HelpListDescription != "" {
			fmt.Printf("	%-13s: %-16s", cmd.Name, cmd.HelpListDescription)

		}

	}

}

func init() {
	RegisterNewCommand(&helplist)

}