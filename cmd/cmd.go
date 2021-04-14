package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/manifoldco/promptui"

	"tssh/store"
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

func addOrSave(body string, args []string, isAdd bool) {
	var err error
	if len(body) == 0 {
		prompt := promptui.Prompt{
			Label:     "please input {user@host}",
			Templates: validateTpl,
			Validate: func(input string) error {
				if strings.Index(input, "@") < 0 {
					return fmt.Errorf("con not decode:%s", input)
				}
				x := strings.Split(input, "@")
				if len(x) != 2 {
					return fmt.Errorf("con not decode:%s", input)
				}
				if len(x[0]) == 0 || len(x[1]) == 0 {
					return fmt.Errorf("user && ip required")
				}
				if address := net.ParseIP(x[1]); address == nil {
					return fmt.Errorf("con not decode ip:%s", x[1])
				}
				return nil
			},
		}
		body, err = prompt.Run()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	//获取添加/覆盖配置需要的参数
	config := ParseArgs(args)
	//检查别名输入情况
	if len(config.Name) == 0 {
		prompt := promptui.Prompt{
			Label:     "please input alias name",
			Templates: validateTpl,
		}
		config.Name, err = prompt.Run()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(config.Name) == 0 {
			config.Name = body
		}
	}
	if isAdd {
		if store.ConfigExists(config.Name) {
			fmt.Println("config", config.Name, "exists")
			return
		}
	}
	//检查密码和密钥输入情况
	if len(config.Pwd) == 0 && len(config.PrivateKeyPath) == 0 {
		prompt := promptui.Prompt{
			Label:     "please input passworde",
			Templates: validateTpl,
			Validate: func(input string) error {
				if len(input) == 0 {
					return fmt.Errorf("pwd required")
				}
				return nil
			},
		}
		config.Pwd, err = prompt.Run()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	user, host := GetUserAndHost(body)
	if len(user) == 0 || len(host) == 0 {
		fmt.Println("user and host required")
		return
	}
	var (
		privateKey []byte
	)
	if len(config.PrivateKeyPath) > 0 {
		privateKey, err = ioutil.ReadFile(config.PrivateKeyPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	cfg := store.NewConfig(config.Name, host, user, config.Pwd, privateKey, config.Port)
	err = store.Set(cfg)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func List() {
	batch, err := store.GetBatchConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	list := batch.List()

	prompt := promptui.Select{
		Label:     "Select config ",
		Items:     list,
		Templates: listTpl,
		Size:      20,
	}
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}
	if err := list[index].Conn(); err != nil {
		fmt.Println(err)
	}
}

func Conn(name string) {
	if len(name) == 0 {
		prompt := promptui.Prompt{
			Label:     "conn config name",
			Templates: validateTpl,
			Validate:  validateFunc,
		}
		var err error
		name, err = prompt.Run()
		if err != nil {
			fmt.Println("get config error", err.Error())
			return
		}
	}
	batch, _ := store.GetBatchConfig()
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
		batch, err := store.GetBatchConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
		list := batch.List()

		prompt := promptui.Select{
			Label:     "Select config ",
			Items:     list,
			Templates: listTpl,
			Size:      20,
		}
		index, _, err := prompt.Run()
		if err != nil {
			fmt.Println("error", err.Error())
			return
		}
		name = list[index].Name
	}
	if err := store.Del(name); err != nil {
		fmt.Println(err)
	}
}
