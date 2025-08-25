package incantations

import (
	"fmt"
	"strings"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

var enchant types.Incantation = types.Incantation{
	Name: "enchant",
	Description: "Enchants an rune",
	RevelioAble: false,
	GrimorieDescription: "",
	Cast: enchantIncantation,

}

func enchantIncantation(grandHall *types.GrandHall, args []string) {
	if len(args) < 3 {
		misc.PanicWarn("Incantation usage: enchant <rune> <value>\n", true)
		return

	}

	rawRunesName := strings.TrimSpace(args[1])
	newVal := strings.TrimSpace(args[2])
	_, exists := grandHall.Chamber.Runes[rawRunesName]
	if !exists {
		misc.PanicWarn(fmt.Sprintf("Rune %s not found, use runes to view avaliable runes in current chamber\n", args[1]), true)
		return

	}

	grandHall.Chamber.SetRuneVal(rawRunesName, newVal)

}

func init() {
	RegisterNewIncantation(&enchant)

}
