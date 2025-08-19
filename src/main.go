package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Yyax13/onTop-C2/src/cmd"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

func main() {
	var mainEnv types.MainEnvType = types.MainEnvType{
		Module: &types.Module{},
		
	}

	toolName, _ := misc.Colors("onTopC2", "red")
	promptSignal, _ := misc.Colors("▶▶", "white_bold")
	prompt := fmt.Sprintf("%v %v ", toolName, promptSignal)

	newErr := misc.ForceClearStdout()
	if newErr != nil {
		misc.PanicWarn(fmt.Sprintf("An error ocurred during stdout forced clear: %v", newErr), false)
		os.Exit(0)

	}

	misc.PrintBanner()
	misc.InitInterruptHandler()
	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		moduleName, _ := misc.Colors(fmt.Sprintf("(%s)", mainEnv.Module.Name), "black")
		if mainEnv.Module.Name != "" {
			prompt = fmt.Sprintf("%v %v %v ", moduleName, toolName, promptSignal)

		}

		fmt.Print(prompt)
		ok := stdinScanner.Scan()
		misc.CtrlDHandler(ok, stdinScanner.Err())

		userInput := strings.Split(strings.TrimSpace(stdinScanner.Text()), " ")
		rawCmd := userInput[0]
		if userInput[0] == "" {
			continue
		
		}

		command, okSec := cmd.AvaliableCommands[rawCmd]
		if !okSec {
			misc.PanicWarn(fmt.Sprintf("Command %s was not found, use help command to view all avaliable commands\n", rawCmd), true)
			continue

		}

		command.Run(&mainEnv, userInput)

	}
}
