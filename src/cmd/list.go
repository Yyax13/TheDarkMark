package cmd

import (
	"fmt"
	"strings"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/modules"
	"github.com/Yyax13/onTop-C2/src/types"
)

var list types.Command = types.Command{
	Name:                "list",
	Description:         "List exploits, modules, sessions and anything who is listable",
	Listable:            false,
	HelpListDescription: "",
	Run:                 listCommand,
}

func listCommand(mainEnv *types.MainEnvType, args []string) {
	if len(args) < 2 {
		misc.PanicWarn("Command correct usage: use <module_name>\n\n", false)
		return

	}

	commandOption := args[1]
	targetCommand, ok := AvaliableCommands[strings.TrimSpace(commandOption)]
	if !ok && !strings.HasPrefix(strings.TrimSpace(strings.ToLower(commandOption)), "session") {
		misc.PanicWarn(fmt.Sprintf("Command %v not found\n\n", commandOption), false)
		return

	}

	var listType string
	var listContent any
	switch {
	case strings.HasPrefix(strings.TrimSpace(strings.ToLower(commandOption)), "session"):
		listType = "BOTs (sessions)"
		listContent = modules.Sessions

	case strings.TrimSpace(strings.ToLower(targetCommand.Name)) == "use":
		listType = "modules"
		listContent = modules.AvaliableModules

	default:
		misc.PanicWarn(fmt.Sprintf("The command %v isn't listable\n\n", targetCommand), false)
		return

	}

	fmt.Printf("Avaliable %v:\n", listType)
	switch t := listContent.(type) {
	case map[string]*types.Module:
		for _, v := range t {
			fmt.Printf("	%-10s %-13s\n", v.Name, v.Description)

		}

		fmt.Print("\n")
		return

	case map[string]*modules.ListenerSession:
		for _, v := range t {
			fmt.Printf("	%-10s %-7s\n", v.ID, v.BotIP)

		}

		fmt.Print("\n")
		return

	default:
		misc.PanicWarn(fmt.Sprintf("The type %s isn't supported for listing, contact the developer\n\n", t), false)
		return

	}

}

func init() {
	RegisterNewCommand(&list)

}
