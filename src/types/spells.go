package types

type SpellMacro struct{
	Macro 					string
	Value					any

}

type SpellMethod struct{
	Name					string
	Description				string
	UsageExample			string
	OperatorSideCommand		string
	ImplantSideCommand		string
}

type Spell struct{
	Name					string
	Description				string
	PayloadAsoluteDirPath	string
	Methods					map[string]*SpellMethod
	InsertCommand			func(ImplantSideCommand string, originalData []byte) ([]byte, error)
	Macros					map[string]*SpellMacro

}

