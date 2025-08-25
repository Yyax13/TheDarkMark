package types

type Rune struct {
	Name        string
	Description string
	Required    bool
	Value       any
}

type Incantation struct {
	Name                	string
	Description         	string
	RevelioAble            bool
	GrimorieDescription 	string
	Cast                 	func(grandHall *GrandHall, args []string)
}
