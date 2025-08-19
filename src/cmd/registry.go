package cmd

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableCommands map[string]*types.Command = make(map[string]*types.Command)

func RegisterNewCommand(cmd *types.Command) {
	AvaliableCommands[cmd.Name] = cmd
	
}
