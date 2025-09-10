package incantations

import (
	"fmt"
	"os"
	"github.com/Yyax13/onTop-C2/src/types"
	"github.com/Yyax13/onTop-C2/src/misc"

)

var scourgify types.Incantation = types.Incantation{
	Name: "scourgify",
	Description: "Clear the stdout and print the banner",
	RevelioAble: false,
	GrimorieDescription: "",
	Cast: scourgifyCommand,

}

func scourgifyCommand(_ *types.GrandHall, _ []string) {
	e := misc.ForceClearStdout()
	if e != nil {
		misc.PanicWarn(fmt.Sprintf("\nAn error ocurred during stdout forced clear: %v", e), true) // this err msg get out of theme scope
		os.Exit(0)

	}

	misc.PrintBanner()

}

func init() {
	RegisterNewIncantation(&scourgify)

}