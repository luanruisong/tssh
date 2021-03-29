package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"tssh/ssh"
	"tssh/store"
)

func main() {
	//ssh.Terminal("172.16.1.156", "root", "Qaz!@#9ol.=[;.", 22)

	args := os.Args
	if len(args) < 2 {
		panic(errors.New("args required"))
	}

	cmd := args[1]

	switch strings.ToLower(cmd) {
	case "add":
		cfg, err := parseCfg(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = store.Add(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("add", cfg.Name, " success")
	case "del":
		name, err := parseName(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = store.Del(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("del", name, " success")
	case "set":
		cfg, err := parseCfg(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = store.Set(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("set", cfg.Name, " success")
	case "get":
		name, err := parseName(args)
		if err != nil {
			fmt.Println(err)
			return
		}
		info, err := store.Get(name)
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

func parseCfg(args []string) (*store.SSHConfig, error) {
	if len(args) < 6 {
		return nil, errors.New("args error")
	}
	name, ip, user, pwd := args[2], args[3], args[4], args[5]
	port := 22
	if len(args) > 6 {
		port, _ = strconv.Atoi(args[6])
	}
	return store.NewConfig(name, ip, user, pwd, port), nil
}

func parseName(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("can not get name")
	}
	return args[2], nil
}
