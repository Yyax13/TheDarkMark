package incantations

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/types"
	"github.com/Yyax13/onTop-C2/src/misc"

)

var cast types.Incantation = types.Incantation{
	Name: "cast",
	Description: "Cast the current chamber",
	RevelioAble: false,
	GrimorieDescription: "",
	Cast: castIncantation,

}

func castIncantation(grandHall *types.GrandHall, args []string) {
	err := grandHall.Chamber.CheckRunesValue()
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("%v\n", err.Error()), true)
		return

	}


	switch grandHall.Chamber.Parallel {
	case true:
		go grandHall.Chamber.Execute(grandHall.Chamber.Runes)
		return

	default:
		grandHall.Chamber.Execute(grandHall.Chamber.Runes)
		return

	}

}

func init() {
	RegisterNewIncantation(&cast)

}
