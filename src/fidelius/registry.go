package fidelius

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableFidelius map[string]*types.Fidelius = make(map[string]*types.Fidelius)
func RegisterNewFidelius(e *types.Fidelius) {
	AvaliableFidelius[e.Name] = e

}

type FideliusCreator func(params map[string]string) (types.FideliusCasting, error)
var AvaliableFideliusCreators map[string]FideliusCreator = make(map[string]FideliusCreator)
func RegisterNewFideliusCreator(name string, creator FideliusCreator) {
	AvaliableFideliusCreators[name] = creator

}
