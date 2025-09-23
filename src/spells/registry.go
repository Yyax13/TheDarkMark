package spells

import "github.com/Yyax13/onTop-C2/src/types"

var AvaliableSpells map[string]*types.Spell = make(map[string]*types.Spell)
func RegisterNewSpell(spell *types.Spell) {
	AvaliableSpells[spell.Name] = spell
	
}
