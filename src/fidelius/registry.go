package fidelius

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableFidelius map[string]*types.Fidelius = make(map[string]*types.Fidelius)
func RegisterNewFidelius(e *types.Fidelius) {
	AvaliableFidelius[e.Name] = e

}