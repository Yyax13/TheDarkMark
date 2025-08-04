package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"os/signal"
	"syscall"

	"github.com/Yyax13/onTop-C2/src/listener"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/session"
)

func main() {
	rawUser, err := user.Current()
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error ocurred when tried to fetch current username: %v", err), false)
		os.Exit(0)

	}

	userName, _ := misc.Colors(rawUser.Username, "red")
	toolName, _ := misc.Colors("onTopC2", "black")
	separator, _ := misc.Colors("@", "white")
	promptSignal, _ := misc.Colors("->", "white")

	newErr := misc.ForceClearStdout()
	if newErr != nil {
		misc.PanicWarn(fmt.Sprintf("An error ocurred during stdout forced clear: %v", newErr), false)
		os.Exit(0)

	}

	misc.PrintBanner()
	misc.InitInterruptHandler()
	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%v%v%v %v ", userName, separator, toolName, promptSignal)
		ok := stdinScanner.Scan()
		misc.CtrlDHandler(ok, stdinScanner.Err())

		userInput := strings.TrimSpace(stdinScanner.Text())

		switch {
		case userInput == "":
			continue

		case strings.HasPrefix(userInput, "use"):
			switch strings.TrimPrefix(userInput, "use ") {
			case "listener":
				misc.SysLog("Insert port: ", false)
				ok := stdinScanner.Scan()
				misc.CtrlDHandler(ok, stdinScanner.Err())

				port, err := strconv.ParseInt(strings.TrimSpace(stdinScanner.Text()), 0, 64)
				if err != nil {
					misc.PanicWarn(fmt.Sprintf("Error: %v\n", err), false)

				}

				if port < 1024 || port > 65535 {
					misc.PanicWarn(fmt.Sprintf("Error: %v", "Can't use port outside the range (1024 ~ 65535)\n"), false)
					continue

				}

				go listener.StartListener(int(port))

			default:
				misc.PanicWarn("Module not found, check avaliable modules using: list use\n", false)
				continue

			}

		case strings.HasPrefix(userInput, "session"):
			botID := strings.TrimSpace(strings.TrimPrefix(userInput, "session "))
			if !strings.HasPrefix(botID, "bot-") {
				misc.SysLog("Usage: session <session_id>\n", false)
				continue

			}

			sessionSigintChannel := make(chan struct{})
			signal.Stop(misc.InterruptSigs)
			signal.Notify(misc.ChanInterruptSigs, syscall.SIGINT)
			session.InteractWithSession(botID, sessionSigintChannel)

		case strings.HasPrefix(userInput, "list"):
			switch strings.TrimPrefix(userInput, "list "){
			case "use":
				avaliableModules := map[string]string{
					"listener": "A TCP Listener with sessions with focus in reverse shells listening",

				}

				fmt.Println("Avaliable modules:")
				for module, description := range avaliableModules {
					fmt.Printf("	%v: %v\n", module, description)

				}

			case "session":
				fmt.Println("Avaliable sessions (IDs):")
				for sessID := range listener.Sessions {
					fmt.Printf("	%v\n", sessID)

				}

			default:
				misc.PanicWarn("Command not found, use helplist to view avaliable commands to list\n", false)

			}
			
		case userInput == "clear":
			e := misc.ForceClearStdout()
			if e != nil {
				misc.PanicWarn(fmt.Sprintf("\nAn error ocurred during stdout forced clear: %v", e), true)
				os.Exit(0)

			}

		misc.PrintBanner()

		case strings.HasPrefix(userInput, "helplist"):
			avaliableCommands := map[string]string{
				"use": "List all modules in use command",
				"session": "List all avaliable sessions IDs",

			}

			fmt.Println("Avaliable commands to list:")
			for command, description := range avaliableCommands {
				fmt.Printf("	%v: %v\n", command, description)

			}

		case strings.HasPrefix(userInput, "help"):
			commands := map[string]string{
				"use": "Use an specifc module",
				"session": "Interact with a bot (session)",
				"list": "List anything avaliable in commands (e.g.: avaliable modules to use command)",
				"helplist": "List commands who're avaliable to be listed",
				"help": "Show this table",

			}

			fmt.Println("Avaliable commands:")
			for command, description := range commands {
				fmt.Printf("	%v: %v\n", command, description)
				
			}

		case userInput == "exit":
			os.Exit(0)

		default:
			misc.PanicWarn("Command not found, use help to view avaliable commands\n", false)
			continue

		}

		fmt.Print("\n")

	}
}
