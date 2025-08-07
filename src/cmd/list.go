package cmd

import (
	"fmt"
	"strings"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/modules"
	"github.com/Yyax13/onTop-C2/src/types"
	
)

var list types.Command = types.Command{
	Name: "list",
	Description: "List exploits, modules, sessions and anything who is listable",
	Listable: false,
	HelpListDescription: "",
	Run: listCommand,

}

func listCommand(mainEnv *types.MainEnvType, args []string) {
	commandOption, ok := opts["COMMAND"] // ADAPTAR P/ ARGS
	if !ok {
		misc.PanicWarn(fmt.Sprintf("Option %s is unset, check it an try again\n", "'COMMAND'"), true)
		return

	}

	targetCommand := commandOption.Value
	if targetCommand == nil {
		misc.PanicWarn(fmt.Sprintf("The '%s' value is unset (nil by default), check it and try again\n", commandOption.Name), true)
		return

	}

	var listType string
	var listContent any;
	switch strings.TrimSpace(strings.ToLower(targetCommand.(string))) {
	case "use":
		listType = "modules"
		listContent = modules.AvaliableModules

	case "session", "sessions":
		listType = "BOTs (sessions)"
		listContent = modules.Sessions

	default:
		misc.PanicWarn(fmt.Sprintf("The command %v isn't listable\n", targetCommand), true)
		return

	}

	fmt.Printf("Avaliable %v:\n", listType)
	switch t := listContent.(type){
	case map[string]types.Module:
		for _, v := range t {
			fmt.Printf("	%-10s %-13s", v.Name, v.Description)

		}

		return

	case map[string]*modules.ListenerSession:
		for _, v := range t {
			fmt.Printf("	%-10s %-7s", v.ID, v.BotIP)
			
		}

		return

	default:
		misc.PanicWarn(fmt.Sprintf("The type %s isn't supported for listing, contact the developer\n", t), true)
		return

	}

}

func init() {
	RegisterNewCommand(&list)

}