package main

/*
Copyright © 2021 Luan Ruisong
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

	"github.com/luanruisong/tssh/cmd"
	"github.com/luanruisong/tssh/store"
)

var version string

func main() {

	//flag.Parse()
	cmd.Logo()
	if len(os.Args) < 2 {
		cmd.Help()
		return
	}

	var (
		flag  = os.Args[1]
		alias string
		args  []string
	)
	if len(os.Args) >= 3 {
		alias = os.Args[2]
		args = os.Args[3:]
	}
	switch flag {
	case "h", "-h", "help", "-help":
		cmd.Help()
	case "v", "-v", "version", "-version":
		fmt.Println("version", version)
	case "e", "-e", "env", "-env":
		store.FmtEnv()
	case "d", "-d", "del", "-del":
		cmd.Del(alias)
	case "a", "-a", "add", "-add":
		cmd.Add(alias, args)
	case "s", "-s", "save", "-save":
		cmd.Save(alias, args)
	case "l", "-l", "list", "-list":
		cmd.List()
	case "c", "-c", "conn", "-conn":
		cmd.Conn(alias)
	default:
		cmd.Help()
	}

}
