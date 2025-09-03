package fidelius

import (
	"github.com/Yyax13/onTop-C2/src/types"

)

type basic_noneFidelius struct{}

func (none basic_noneFidelius) Encode(data []byte) ([]byte, error) {
	return data, nil

}

func (none basic_noneFidelius) Decode(data []byte) ([]byte, error) {
	return data, nil

}

var Basic_none types.Fidelius = types.Fidelius{
	Name: "basic/none",
	Description: "Don't really encode nothing, just return the original payload",
	Fidelius: basic_noneFidelius{},

}

func basic_noneCreator() (FideliusCreator) {
	return func(params map[string]string) (types.FideliusCasting, error) {
		return &basic_noneFidelius{}, nil

	}

}

func init() {
	RegisterNewFidelius(&Basic_none)
	RegisterNewFideliusCreator("basic/none", basic_noneCreator())
	
}