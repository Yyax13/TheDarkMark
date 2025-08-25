package incantations

import (
	"fmt"
	"github.com/Yyax13/onTop-C2/src/types"
)

var grimorie types.Incantation = types.Incantation{
	Name:                "grimorie",
	Description:         "Shows all incantations compatible with revelio incantation",
	RevelioAble:            false,
	GrimorieDescription: "",
	Cast:                 grimorieIncantation,
}

func grimorieIncantation(_ *types.GrandHall, _ []string) {
	fmt.Println("Revelio able incantations:")
	for _, incantation := range AvaliableIncantations {
		if incantation.RevelioAble && incantation.GrimorieDescription != "" {
			fmt.Printf("	%-13s: %-16s\n", incantation.Name, incantation.GrimorieDescription)

		}

	}

	fmt.Print("\n")

}

func init() {
	RegisterNewIncantation(&grimorie)

}
