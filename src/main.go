package main

import (
	"fmt"
	"os"
	"strings"
	"io"

	"github.com/Yyax13/onTop-C2/src/incantations"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"

	"github.com/chzyer/readline"
)

func main() {
	var grandHall types.GrandHall = types.GrandHall{
		Chamber: &types.Chamber{},
	}

	toolName, _ := misc.Colors("Mark", "green")
	promptSignal, _ := misc.Colors("⌁", "dark_green_bold")
	prompt := fmt.Sprintf("%v %v %v", toolName, promptSignal, misc.AvaliableColors["white"])

	newErr := misc.ForceClearStdout()
	if newErr != nil {
		misc.PanicWarn(fmt.Sprintf("An error ocurred during scourgify: %v", newErr), false)
		os.Exit(0)

	}

	misc.PrintBanner()
	misc.InitInterruptHandler()

	rl, ee := readline.New(prompt) // readline broke parallel stdout cause print the prompt every time and parallel can't print the stdout (fr it can, but it's ugly)
	if ee != nil {
		fmt.Println("Some error occurred during readline initialization: ", ee)
		os.Exit(0)

	}

	defer func() { _ = rl.Close() }()
	for {
		chamberName, _ := misc.Colors(fmt.Sprintf("(%s)", grandHall.Chamber.Name), "dark_green_bold")
		if grandHall.Chamber.Name != "" {
			rl.SetPrompt(fmt.Sprintf("%v %v %v ", chamberName, toolName, promptSignal))

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

		command, okSec := incantations.AvaliableIncantations[rawCmd]
		if !okSec {
			misc.PanicWarn(fmt.Sprintf("Incantation %s was not found, use pensieve incantation to view all avaliable incantations\n\n", rawCmd), false)
			continue

		}

		command.Cast(&grandHall, userInput)

	}

}
