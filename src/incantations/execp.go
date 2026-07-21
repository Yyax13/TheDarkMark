package incantations

import (
    "fmt"
    "os"
    "path/filepath"
    "bufio"
    "strings"
    "runtime"

    "github.com/Yyax13/onTop-C2/src/types"
    "github.com/Yyax13/onTop-C2/src/misc"
)

var execp types.Incantation = types.Incantation{
    Name: "execp",
    Description: "Execute commands from a profile (at ./profiles)",
    RevelioAble: false,
    GrimorieDescription: "",
    Cast: execpIncantation,
}

func execpIncantation(grandHall *types.GrandHall, args []string) {
    if len(args) < 2 {
        misc.PanicWarn("Incantation correct usage: execp <profile_name>\n\n", false)
        return
    }

    _, _filename, _, _ := runtime.Caller(0)
	_dirname := filepath.Dir(_filename)

    profilePath := filepath.Join(_dirname, "..", "..", "profiles", args[1])
    file, err := os.Open(profilePath)
    if err != nil {
        if os.IsNotExist(err) {
            misc.PanicWarn(fmt.Sprintf("The profile %s don't exist", args[1]), true)
            return
        }

        misc.PanicWarn(err.Error(), true)
        return
    }

    defer file.Close()
    scanner := bufio.NewScanner(file)

    misc.SysLog(fmt.Sprintf("Running profile %s\n", args[1]), true)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.TrimSpace(line) == "" {
            continue
        }

        lInput := strings.Split(strings.TrimSpace(line), " ")
        rawCmd := lInput[0]
        if rawCmd == "" {
            continue
        }

        if rawCmd == "execp" {
            misc.PanicWarn("Can't run recursive execp, skipping", true)
            continue
        }

        cmd, ok := AvaliableIncantations[rawCmd]
        if !ok {
			misc.PanicWarn(fmt.Sprintf("Incantation %s was not found, skipping", rawCmd), true)
            continue
        }

        cmd.Cast(grandHall, lInput)
    }
}

func init() {
    RegisterNewIncantation(&execp)
}
