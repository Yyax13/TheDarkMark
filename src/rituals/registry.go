package rituals

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableRituals map[string]*types.Ritual = make(map[string]*types.Ritual)
func RegisterNewRitual(r *types.Ritual) {
	AvaliableRituals[r.Name] = r
	
}