package misc

import "fmt"

func PanicWarn(logContent string, newLine bool) {
	var warnText, _ = Colors("[WARN]", "red")
	if newLine {
		warnText = "\n" + warnText

	}

	logContentColored, _ := Colors(logContent, "yellow")
	fmt.Printf("%s %v", warnText, logContentColored)

}

func SysLog(logContent string, newLine bool) {
	var sysText, _ = Colors("[SYSTEM]", "black")
	if newLine {
		sysText = "\n" + sysText

	}
	
	logContentColored, _ := Colors(logContent, "white")

	fmt.Printf("%s %v", sysText, logContentColored)

}
