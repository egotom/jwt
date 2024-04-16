package models

type Config struct{
	DSN					string 					`yaml:"DSN"`
	Secret 				string					`yaml:"Secret"` 
	APIKey 				map[string] string		`yaml:"APIKey"` 
	GinMode 			string 					`yaml:"GinMode"` 
	Port 				int 					`yaml:"Port"` 
	HTTPProxy 			string 					`yaml:"HTTPProxy"` 
	Greets 				[]string				`yaml:"Greets"` 
	Prompts 			map[string][]string		`yaml:"Prompts"` 
	EnableSamllTalk		bool					`yaml:"EnableSamllTalk"` 
	EnableInvite		bool					`yaml:"EnableInvite"` 
}

