package misc

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	
)

var InterruptSigs chan os.Signal = make(chan os.Signal, 1)

func InitInterruptHandler() {
	signal.Notify(InterruptSigs, syscall.SIGINT)

	go func() {
		<-InterruptSigs
		SysLog("Quiting...\n", true)
		os.Exit(0)

	}()

}

var ChanInterruptSigs chan os.Signal = make(chan os.Signal, 1)

func ChanInterruptHandler(channel chan struct{}) {
	signal.Notify(ChanInterruptSigs, syscall.SIGINT)

	go func() {
		<-ChanInterruptSigs
		fmt.Print("\n")
		channel <- struct{}{}

	}()

}

func CtrlDHandler(o bool, e error) {
	if !o {
		if e == nil {
			fmt.Print("\n")
			os.Exit(0)

		}

		os.Exit(1)

	}

}

func ChanCtrlDHandler(o bool, e error, channel chan struct{}) {
	if !o {
		if e == nil {
			fmt.Print("\n")
			channel <- struct{}{}

		}

		PanicWarn(fmt.Sprintf("An error occurred during attacker input: %v", e), true) // need to change from attacker to caster or smt like that
		channel <- struct{}{}

	}

}