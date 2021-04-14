package cmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"

	"tssh/store"
)

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

  env		get env info 				(e|-e)
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
	detail = `
----------------------------------------------------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Ip:" | faint }}	{{ .Ip }}
{{ "User:" | faint }}	{{ .User }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "ConnMode:" | faint }}	{{ .ConnMode }}
{{ "SaveAt:" | faint }}	{{ .SaveAt }}`
)

func Logo() {
	fmt.Print(logoStr)
}

func Help() {
	fmt.Print(helpStr)
}

var (
	validateFunc = func(input string) error {
		g, err := store.GetBatchConfig()
		if err != nil {
			return err
		}
		if g.Get(input) == nil {
			return errors.New("can not get config")
		}
		return err
	}

	validateTpl = &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}
)
