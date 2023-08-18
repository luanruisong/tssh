package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/luanruisong/tssh/constant"
	"github.com/luanruisong/tssh/store"
)

type (
	CmdSSHConfig struct {
		Name           string
		Port           int
		Pwd            string
		PrivateKeyPath string
	}
)

func ParseArgs(args []string) *CmdSSHConfig {
	var (
		fs = flag.NewFlagSet("set", flag.ExitOnError)
		n  = fs.String("n", "", "set name in (-a|-s)")
		p  = fs.String("p", "", "set password in (-a|-s)")
		P  = fs.Int("P", 22, "set port in (-a|-s)")
		k  = fs.String("k", "", "set private_key path in (-a|-s)")
	)
	fs.Parse(args)
	return &CmdSSHConfig{
		Name:           *n,
		Port:           *P,
		Pwd:            *p,
		PrivateKeyPath: *k,
	}
}

func GetUserAndHost(a string) (string, string) {
	if len(a) > 0 {
		if idx := strings.Index(a, "@"); idx > 0 {
			return a[:idx], a[idx+1:]
		}
	}
	return "", ""
}

func Add(body string, args []string) {
	addOrSave(body, args, true)
}

func Save(body string, args []string) {
	addOrSave(body, args, false)
}

func List() {
	list := store.ListConfig()
	if len(list) == 0 {
		fmt.Println("can not get config list")
		return
	}
	prompt := promptui.Select{
		Label:     "Connect config ",
		Items:     list,
		Templates: constant.ListTpl,
		Size:      20,
	}
	_, name, err := prompt.Run()
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	if conn := store.GetConfig(name); conn != nil {
		if err := conn.Conn(); err != nil {
			fmt.Println(err)
		}
	}
}

func Conn(name string) {
	if len(name) == 0 {
		List()
		return
	}
	batch := store.GetBatchConfig()
	info := batch.Get(name)
	if info == nil {
		fmt.Println("can not find config", name)
		return
	}
	if err := info.Conn(); err != nil {
		fmt.Println(err)
	}
}

func Del(name string) {
	if len(name) == 0 {
		list := store.ListConfig()
		prompt := promptui.Select{
			Label:     "Delete config ",
			Items:     list,
			Templates: constant.ListTpl,
			Size:      20,
		}
		var err error
		_, name, err = prompt.Run()
		if err != nil {
			fmt.Println("error", err.Error())
			return
		}
	}
	if err := store.Del(name); err != nil {
		fmt.Println(err)
	}
}
