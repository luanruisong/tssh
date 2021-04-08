package cmd

import "fmt"

const (
	logoStr = `
 ______   ______     ______     __  __    
/\__  _\ /\  ___\   /\  ___\   /\ \_\ \   
\/_/\ \/ \ \___  \  \ \___  \  \ \  __ \  
   \ \_\  \/\_____\  \/\_____\  \ \_\ \_\ 
    \/_/   \/_____/   \/_____/   \/_/\/_/

`
	helpStr = `
Usage of TSSH:

  env		get evn info 				(e|-e)
  version	get version info			(v|-v)
  list 		get config list				(l|-l)
  conn		connect to alias			(c|-c)
  delete 	del config by alias			(d|-d)
  add 		add config {user@host}			(a|-a)
  save 		reset config {user@host}		(s|-s)
	  -P int
			set port in (add|save) (default 22)
	  -k string
			set private_key path in (add|save)
	  -n string
			set alias name in (add|save)
	  -p string
			set password in (add|save)

`
)

func Logo() {
	fmt.Print(logoStr)
}

func Help() {
	fmt.Print(helpStr)
}

func LogoAndHelp() {
	Logo()
	Help()
}
