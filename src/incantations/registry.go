package incantations

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableIncantations map[string]*types.Incantation = make(map[string]*types.Incantation)

func RegisterNewIncantation(cmd *types.Incantation) {
	AvaliableIncantations[cmd.Name] = cmd
	
}
