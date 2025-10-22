package incantations

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Yyax13/onTop-C2/src/chambers"
	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/rituals"
	"github.com/Yyax13/onTop-C2/src/spells"
	"github.com/Yyax13/onTop-C2/src/types"
)

var revelio types.Incantation = types.Incantation{
	Name:                "revelio",
	Description:         "List spells, chambers, imperius and anything who is revelio able",
	RevelioAble:            false,
	GrimorieDescription: "",
	Cast:                 revelioIncantation,
}

func revelioIncantation(_ *types.GrandHall, args []string) {
	if len(args) < 2 {
		misc.PanicWarn("Incantation correct usage: revelio <incantation_or_chamber_name>\n\n", false)
		return

	}

	incantationRunes := args[1]
	inputFromUser := strings.TrimSpace(strings.ToLower(incantationRunes))
	targetIncantation, ok := AvaliableIncantations[strings.TrimSpace(incantationRunes)]
	if !ok && !(
			strings.HasPrefix(inputFromUser, "imperius") 	|| 
			strings.HasPrefix(inputFromUser, "inferi") 		|| 
			strings.HasPrefix(inputFromUser, "fidelius") 	|| 
			strings.HasPrefix(inputFromUser, "ritual")		||
			strings.HasPrefix(inputFromUser, "spell")) {
		misc.PanicWarn(fmt.Sprintf("Incantation %v not found\n\n", incantationRunes), false)
		return

	}

	var revelioType string
	var revelioContent any
	switch { // brooo i hate switchs but idk other way for doing that
	case strings.HasPrefix(inputFromUser, "imperius"), strings.HasPrefix(inputFromUser, "inferi"):
		revelioType = "Inferis (imperius)"
		revelioContent = chambers.Inferis

	case strings.HasPrefix(inputFromUser, "wield"):
		revelioType = "Chambers"
		revelioContent = chambers.AvaliableModules

	case strings.HasPrefix(inputFromUser, "fidelius"):
		revelioType = "Fidelius"
		revelioContent = fidelius.AvaliableFidelius

	case strings.HasPrefix(inputFromUser, "ritual"):
		revelioType = "Rituals"
		revelioContent = rituals.AvaliableRituals
	
	case strings.HasPrefix(inputFromUser, "spell"):
		revelioType = "Spells"
		revelioContent = spells.AvaliableSpells

	default:
		misc.PanicWarn(fmt.Sprintf("The incantation %v isn't revelio able\n\n", targetIncantation), false)
		return

	}

	fmt.Printf("Avaliable %v:\n", revelioType)
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	switch t := revelioContent.(type) {
	case map[string]*types.Chamber:
		for _, v := range t {
			fmt.Fprintf(writer, "	\t%s\t%s\n", v.Name, v.Description)

		}

	case map[string]*chambers.Inferi:
		fmt.Fprintf(writer, "	\t%s\t%s\n", "INFERI ID", "IP")
		for _, v := range t {
			fmt.Fprintf(writer, "	\t%s\t%s\n", v.ID, v.BotIP)

		}

	case map[string]*types.Fidelius:
		fmt.Fprintf(writer, "	\t%s\t%s\n", "Name", "Description")
		for _, v := range t {
			fmt.Fprintf(writer, "	\t%s\t%s\n", v.Name, v.Description)
			
		}

	case map[string]*types.Ritual:
		fmt.Fprintln(writer, "	\tName\tDescription\tDefault Fidelius")
		for _, v := range t {
			fmt.Fprintf(writer, "	\t%s\t%s\t%s\n", v.Name, v.Description, v.Encoder.Name)
			
		}
	
	case map[string]*types.Spell:
		fmt.Fprintf(writer, "	\tName\tDescription\n")
		for _, v := range t {
			fmt.Fprintf(writer, "	\t%s\t%s\n", v.Name, v.Description)

		}

	default:
		misc.PanicWarn(fmt.Sprintf("The type %s isn't supported for Revelio, contact the witch (open a issue)\n\n", t), false)
		return

	}

	writer.Flush()
	fmt.Print("\n")

}

func init() {
	RegisterNewIncantation(&revelio)

}
