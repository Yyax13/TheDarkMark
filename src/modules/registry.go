package modules

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableModules map[string]*types.Module = make(map[string]*types.Module)

func RegisterNewModule(mdl *types.Module) {
	AvaliableModules[mdl.Name] = mdl

}
