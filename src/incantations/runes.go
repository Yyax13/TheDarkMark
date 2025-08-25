package incantations

import (
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

var runes types.Incantation = types.Incantation{
	Name:                "runes",
	Description:         "Show the runes in current chamber",
	RevelioAble:            false,
	GrimorieDescription: "",
	Cast:                 runesIncantation,
}

func runesIncantation(grandHall *types.GrandHall, args []string) {
	if grandHall.Chamber.Name == "" {
		misc.PanicWarn("You aren't in a chamber, run revelio cast to view avaliable chambers.\n", true)
		return

	}

	grandHall.Chamber.ListAvaliableRunes()

}

func init() {
	RegisterNewIncantation(&runes)

}
