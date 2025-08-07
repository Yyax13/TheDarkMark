package cmd

import (
	"fmt"
	"strings"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/modules"
	"github.com/Yyax13/onTop-C2/src/types"

)

var use types.Command = types.Command{
	Name: "use",
	Description: "Interact to a specifc module",
	Listable: true,
	HelpListDescription: "List all modules who're usable by use command",
	Run: useCommand,

}

func useCommand(mainEnv *types.MainEnvType, cmdArgs []string) {
	if len(cmdArgs) < 2 {
		misc.PanicWarn(fmt.Sprintf("Command correct usage: use <module_name>\n"), true)
		return

	}

	targetModuleByUser := cmdArgs[1]
	targetModule, ok := modules.AvaliableModules[strings.TrimSpace(strings.ToLower(targetModuleByUser))]
	if !ok {
		misc.PanicWarn(fmt.Sprintf("Not found the module %v, check if it exists using: list use\n", targetModuleByUser), true)
		return

	}



}