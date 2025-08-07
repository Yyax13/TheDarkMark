package types

type Option struct{
	Name 			string
	Description 	string
	Required 		bool
	Value 			any

}

type Command struct{
	Name 					string
	Description 			string
	Listable 				bool
	HelpListDescription		string
	Run 					func(mainEnv *MainEnvType, args []string)

}

