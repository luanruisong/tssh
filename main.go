package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"tssh/ssh"
	"tssh/store"
)

func main() {

	args := os.Args
	if len(args) < 2 {
		panic(errors.New("args required"))
	}

	cmd := args[1]

	switch strings.ToLower(cmd) {
	case "env":
		store.Env()
	case "add":
		err := store.Add(args)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "del":
		err := store.Del(args)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "set":
		err := store.Set(args)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "get":
		info, err := store.GetByArgs(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		printCfg([]store.SSHConfig{*info})
	case "list":
		list, err := store.List()
		if err != nil {
			fmt.Println(err)
			return
		}
		printCfg(list)
	default:
		if len(args) < 2 {
			fmt.Println(fmt.Errorf("can not get name"))
			return
		}
		info, err := store.Get(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		ssh.Terminal(info.Ip, info.User, info.Pwd, info.Port)
	}

}

func printCfg(cfgs []store.SSHConfig) {
	w := tabwriter.NewWriter(os.Stdout, 10, 3, 5, ' ',
		tabwriter.AlignRight)
	fmt.Fprintln(w, "No\tname\tip\tuser\tpwd\tport\tsave_at\t")
	for i, v := range cfgs {
		s := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d\t%s\t", i+1, v.Name, v.Ip, v.User, v.Pwd, v.Port, v.SaveAt)
		fmt.Fprintln(w, s)
	}
	w.Flush()
}
