package fidelius

import (
	"fmt"
	"crypto/rand"
	"crypto/sha512"
	"io"

	"github.com/Yyax13/onTop-C2/src/misc"
	"github.com/Yyax13/onTop-C2/src/types"
)

type basic_xorFidelius struct{}

func (xe basic_xorFidelius) Encode(data []byte) ([]byte, error) {
	nonce := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return data, fmt.Errorf("some error occurred in basic/xor encode")

	}

	h := sha512.New()
	h.Write(nonce)
	h.Write(fmt.Appendf([]byte{}, "len:%d", len(data)))
	dynamicKey := h.Sum(nil)
	xorData := misc.Xor(data, dynamicKey)
	return append(nonce, xorData...), nil

}

func (xe basic_xorFidelius) Decode(data []byte) ([]byte, error) {
	nonceSize := 32
	if len(data) < nonceSize {
		return data, fmt.Errorf("illegal payload: payload is smaller than nonce")

	}

	nonce := data[:nonceSize]
	xorData := data[nonceSize:]
	h := sha512.New()
	h.Write(nonce)
	h.Write(fmt.Appendf([]byte{}, "len:%d", len(xorData)))
	dynamicKey := h.Sum(nil)
	plain := misc.Xor(xorData, dynamicKey)

	return plain, nil

}

var Basic_xor types.Fidelius = types.Fidelius{
	Name: "basic/xor",
	Description: "Just xor the strings with dynamic key",
	Fidelius: basic_xorFidelius{},

}

func basic_xorCreator() (FideliusCreator) {
	return func(params map[string]string) (types.FideliusCasting, error) {
		return &basic_xorFidelius{}, nil
	
	}

}

func init() {
	RegisterNewFidelius(&Basic_xor)
	RegisterNewFideliusCreator("basic/xor", basic_xorCreator())

}