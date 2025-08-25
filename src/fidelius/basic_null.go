package fidelius

import (
	"github.com/Yyax13/onTop-C2/src/types"

)

type basic_nullFidelius struct{}

func (null basic_nullFidelius) Encode(data []byte) ([]byte, error) {
	result := []byte{}
	for _, b := range data {
		result = append(result, 0x01)
		for range int(b) {
			result = append(result, 0x00)

		}
		result = append(result, 0x02)
		
	}
	return result, nil

}

func (null basic_nullFidelius) Decode(data []byte) ([]byte, error) {
	result := []byte{}
	var i uint8
	mu := false

	for _, b := range data {
		switch b {
		case 0x00:
			if mu {
				i++
			
			}

		case 0x01:
			i = 0
			mu = true
		
		case 0x02:
			if mu {
				result = append(result, i)
				mu = false

			}
			
		}

	}
	return result, nil

}

var Basic_null types.Fidelius = types.Fidelius{
	Name: "basic/null",
	Description: "Changes everything to \"null\" (0x00)",
	Fidelius: basic_nullFidelius{},

}

func init() {
	RegisterNewFidelius(&Basic_null)
	
}