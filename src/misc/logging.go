package misc

import "fmt"

func PanicWarn(logContent string, newLine bool) {
	if newLine {
		fmt.Println()

	}

	logContentColored, _ := Colors(logContent, "yellow_bold")
	fmt.Printf("%v", logContentColored)

}

func SysLog(logContent string, newLine bool) {
	if newLine {
		fmt.Println()

	}
	
	logContentColored, _ := Colors(logContent, "green")

	fmt.Printf("%v", logContentColored)

}
