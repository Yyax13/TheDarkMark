package types

type FideliusCasting interface{
	Encode(payload []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)

}

type Fidelius struct{
	Name				string
	Description			string
	Fidelius  			FideliusCasting

}