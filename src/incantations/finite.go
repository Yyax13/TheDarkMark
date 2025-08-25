package incantations

import (
	"os"
	"github.com/Yyax13/onTop-C2/src/types"

)

var finite types.Incantation = types.Incantation{
	Name: "finite",
	Description: "Exit the program without any messages",
	RevelioAble: false,
	GrimorieDescription: "",
	Cast: finiteIncantation,
	
}

func finiteIncantation(_ *types.GrandHall, _ []string) {
	os.Exit(0)

}

func init() {
	RegisterNewIncantation(&finite)

}