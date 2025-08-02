package misc

import (
	"fmt"
    "os/exec"
	"runtime"
	"os"

)

func ClearStdout() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	
}

func ForceClearStdout() error {
	var clearCommand *exec.Cmd

	switch runtime.GOOS{
	case "windows":
		clearCommand = exec.Command("cls")

	default:
		clearCommand = exec.Command("clear")

	}

	clearCommand.Stdout = os.Stdout
	clearCommand.Stderr = os.Stderr

	return clearCommand.Run()

}

func PrintBanner() {
	rawBanner := `
                     _/_/_/_/_/                         _/_/_/    _/_/   
    _/_/    _/_/_/      _/      _/_/    _/_/_/       _/        _/    _/  
 _/    _/  _/    _/    _/    _/    _/  _/    _/     _/            _/     
_/    _/  _/    _/    _/    _/    _/  _/    _/     _/          _/        
 _/_/    _/    _/    _/      _/_/    _/_/_/         _/_/_/  _/_/_/_/     
                                    _/                                   
                                   _/                                    
						`
	
	byText, _ := Colors("by", "black")	
	ownerName, _ := Colors("hoWo the Lammer", "yellow")
	infoText := fmt.Sprintf("\x1b[3m%s \x1b[3m%s\x1b[0m %s\n", byText, ownerName, avaliableColors["reset"])
	banner, _ := Colors(fmt.Sprintf("%s%s", rawBanner, infoText), "red")
	fmt.Printf("%s\n", banner)

}
