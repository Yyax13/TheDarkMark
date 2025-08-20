package main

import (
	"fmt"
	"os"
	"strings"
	"io"

	"github.com/Yyax13/onTop-C2/src/cmd"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"

	"github.com/chzyer/readline"
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
	
	rl, ee := readline.New(prompt)
	if ee != nil {
		fmt.Println("Some error occurred during readline initialization: ", ee)
		os.Exit(0)

	}

	defer func() { _ = rl.Close() }()

	for {
		moduleName, _ := misc.Colors(fmt.Sprintf("(%s)", mainEnv.Module.Name), "black")
		if mainEnv.Module.Name != "" {
			rl.SetPrompt(fmt.Sprintf("%v %v %v ", moduleName, toolName, promptSignal))

		} else {
			rl.SetPrompt(prompt)

		}

		l, err := rl.Readline()
		switch err {
		case io.EOF:
			return
		
		case readline.ErrInterrupt:
			continue
		
		}

		userInput := strings.Split(strings.TrimSpace(l), " ")
		rawCmd := userInput[0]
		if userInput[0] == "" {
			continue

		}

		command, okSec := cmd.AvaliableCommands[rawCmd]
		if !okSec {
			misc.PanicWarn(fmt.Sprintf("Command %s was not found, use help command to view all avaliable commands\n", rawCmd), false)
			continue

		}

		command.Run(&mainEnv, userInput)

	}

}
