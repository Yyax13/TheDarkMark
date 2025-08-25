package chambers

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableModules map[string]*types.Chamber = make(map[string]*types.Chamber)

func RegisterNewModule(mdl *types.Chamber) {
	AvaliableModules[mdl.Name] = mdl

}
