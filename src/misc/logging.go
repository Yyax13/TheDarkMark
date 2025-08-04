package misc

import "fmt"

func PanicWarn(logContent string, newLine bool) {
	if newLine {
		fmt.Println()

	}

	logContentColored, _ := Colors(logContent, "yellow")
	fmt.Printf("%v", logContentColored)

}

func SysLog(logContent string, newLine bool) {
	if newLine {
		fmt.Println()

	}
	
	logContentColored, _ := Colors(logContent, "white")

	fmt.Printf("%v", logContentColored)

}
