package types

import (
	"fmt"
	"net"
)

type FlooNetwork interface{
	Send(data []byte) (error)
	Receive() ([]byte, error)
	Close() (error)
	IsActive() (bool)

}

type Scroll	struct{ // i'll have a hard work to implement a func that fetch all that data from OS in spells
	IP			net.IP
	CPU			struct{
		Name		string
		Cores		int
		Threads		int
		Arch		string
		Clock		int
		Cache		int

	}

	OS			struct{
		Name		string
		Version		string
		Arch		string
		Hostname	string
		Username	string
		Domain		string
		Uptime		int
		AV	struct{
			Name		string
			Active		bool

		}

	}

}

type ArcaneLink struct{
	Network				FlooNetwork
	ClientScroll		*Scroll
	Fidelius			FideliusCasting

}

// need to remove stdout for stealth in spells in every method
func (a *ArcaneLink) Send(data []byte) (error) {
	encodedData, err := a.Fidelius.Encode(data)
	if err != nil {
		return fmt.Errorf("arcanelink: encode failed: %v", err)

	}

	return a.Network.Send(encodedData)

}

func (a *ArcaneLink) Receive() ([]byte, error) {
	received, err := a.Network.Receive()
	if err != nil {
		return nil, fmt.Errorf("arcanelink: failed during data receive: %v", err)

	}

	decodedData, secErr := a.Fidelius.Decode(received)
	if secErr != nil {
		return nil, fmt.Errorf("arcanelink: failed in received data decode: %v", err)

	}

	return decodedData, nil

}

func (a *ArcaneLink) Close() (error) {
	return a.Network.Close()

}

func (a *ArcaneLink) IsActive() bool {
	return a.Network.IsActive()

}

func (a *ArcaneLink) GetScroll() Scroll {
	return *a.ClientScroll

}

func (a *ArcaneLink) SetScroll(newScroll *Scroll) error {
	if newScroll == nil {
		return fmt.Errorf("can't use a nil scroll as clientscroll")
		
	}

	a.ClientScroll = newScroll
	return nil

}

type RitualInit interface{
	InitArcane() (*ArcaneLink, error)

}

type Ritual struct{
	Name	 			string
	Description 		string
	Init				RitualInit
	Fidelius			Fidelius

}
