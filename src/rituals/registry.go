package rituals

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableRituals map[string]*types.Ritual = make(map[string]*types.Ritual)
func RegisterNewRitual(r *types.Ritual) {
	AvaliableRituals[r.Name] = r
	
}

type RitualCreator func(params map[string]string) (types.RitualInit, types.RitualListener, error)
var AvaliableRitualCreators map[string]RitualCreator = make(map[string]RitualCreator)
func RegisterNewRitualCreator(name string, creator RitualCreator) {
	AvaliableRitualCreators[name] = creator

}