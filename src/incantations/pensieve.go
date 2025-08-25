package incantations

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/types"
)

var pensieve types.Incantation = types.Incantation{
	Name:                "pensieve",
	Description:         "Shows all avaliable incantations",
	RevelioAble:            false,
	GrimorieDescription: "",
	Cast:                 pensieveCommand,
}

func pensieveCommand(_ *types.GrandHall, _ []string) {
	fmt.Println("Avaliable incantations:")
	for _, incantations := range AvaliableIncantations {
		fmt.Printf("	%-13s %-16s\n", fmt.Sprintf("%s:", incantations.Name), incantations.Description)

	}

	fmt.Print("\n")

}

func init() {
	RegisterNewIncantation(&pensieve)

}
