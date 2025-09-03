package fidelius

import (
	"github.com/Yyax13/onTop-C2/src/types"

)

type basic_bjumpFidelius struct{}

func (bj basic_bjumpFidelius) Encode(data []byte) ([]byte, error) {
	result := make([]byte, len(data))
	for i := range len(data) {
		result[i] = byte(data[i] + byte(i%256))

	}

	return result, nil

}

func (bj basic_bjumpFidelius) Decode(data []byte) ([]byte, error) {
	result := make([]byte, len(data))
	for i := range len(data) {
		result[i] = byte(data[i] - byte(i%256))

	}

	return result, nil

}

var Basic_bjump types.Fidelius = types.Fidelius{
	Name: "basic/bjump",
	Description: "Just dynamic jumps bytes",
	Fidelius: basic_bjumpFidelius{},
}

func basic_bjumpCreator() (FideliusCreator) {
	return func(params map[string]string) (types.FideliusCasting, error) {
		return &basic_bjumpFidelius{}, nil

	}

}

func init() {
	RegisterNewFidelius(&Basic_bjump)
	RegisterNewFideliusCreator("basic/bjump", basic_bjumpCreator())
	
}