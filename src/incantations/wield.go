package incantations

import (
	"fmt"
	"strings"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/chambers"
	"github.com/Yyax13/onTop-C2/src/types"

)

var wield types.Incantation = types.Incantation{
	Name: "wield",
	Description: "Interact to a specifc module",
	RevelioAble: true,
	GrimorieDescription: "List all chambers who're usable by wield incantations",
	Cast: wieldIncantation,

}

func wieldIncantation(grandHall *types.GrandHall, cmdArgs []string) {
	if len(cmdArgs) < 2 {
		misc.PanicWarn("Incantation correct usage: wield <chamber_name>\n\n", false)
		return

	}

	targetModuleByUser := cmdArgs[1]
	targetModule, ok := chambers.AvaliableModules[strings.TrimSpace(strings.ToLower(targetModuleByUser))]
	if !ok {
		misc.PanicWarn(fmt.Sprintf("Not found the chamber %v, check if it exists using: revelio wield\n\n", targetModuleByUser), false)
		return

	}

	grandHall.Chamber = targetModule

}

func init() {
	RegisterNewIncantation(&wield)

}