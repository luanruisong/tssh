package main

/*
Copyright © 2020 Luan Ruisong
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import (
	"fmt"
	"os"

	"tssh/cmd"
	"tssh/store"
)

var version string

const help = `
Usage of TSSH:

 ______   ______     ______     __  __    
/\__  _\ /\  ___\   /\  ___\   /\ \_\ \   
\/_/\ \/ \ \___  \  \ \___  \  \ \  __ \  
   \ \_\  \/\_____\  \/\_____\  \ \_\ \_\ 
    \/_/   \/_____/   \/_____/   \/_/\/_/

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

func main() {

	//flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(help)
		return
	}

	flag := os.Args[1]
	var alias string
	if len(os.Args) >= 3 {
		alias = os.Args[2]
	}
	switch flag {
	case "h", "-h", "help", "-help":
		fmt.Println(help)
	case "v", "-v", "version", "-version":
		fmt.Println("version", version)
	case "e", "-e", "env", "-env":
		store.Env()
	case "d", "-d", "del", "-del":
		err := store.Del(alias)
		if err != nil {
			fmt.Println(err)
		}
	case "a", "-a", "add", "-add":
		cmd.Add(alias, os.Args[3:])
	case "s", "-s", "save", "-save":
		cmd.Save(alias, os.Args[3:])
	case "l", "-l", "list", "-list":
		cmd.List()
	case "c", "-c", "conn", "-conn":
		cmd.Conn(alias)
	}

	//检测环境变量
	if err := store.DefaultCheck(); err != nil {
		fmt.Println(err)
		return
	}
}
