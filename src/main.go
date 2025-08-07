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
	toolName, _ := misc.Colors("onTopC2", "black")
	promptSignal, _ := misc.Colors("->", "white")
	prompt := fmt.Sprintf("%v %v ", toolName, promptSignal)

	newErr := misc.ForceClearStdout()
	if newErr != nil {
		misc.PanicWarn(fmt.Sprintf("An error ocurred during stdout forced clear: %v", newErr), false)
		os.Exit(0)

	}

	misc.PrintBanner()
	misc.InitInterruptHandler()
	stdinScanner := bufio.NewScanner(os.Stdin)
	var mainEnv types.MainEnvType = types.MainEnvType{
		Module: nil,
		
	}

	for {
		fmt.Print(prompt)
		ok := stdinScanner.Scan()
		misc.CtrlDHandler(ok, stdinScanner.Err())

		userInput := strings.Split(strings.TrimSpace(stdinScanner.Text()), " ")


	}
}
