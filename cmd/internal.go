package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/manifoldco/promptui"

	"tssh/store"
)

var (
	validateFunc = func(input string) error {
		g := store.GetBatchConfig()
		if g.Get(input) == nil {
			return errors.New("can not get config")
		}
		return nil
	}

	validateTpl = &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}
)

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

	if err = store.NewConfig(config.Name, host, user, config.Pwd, privateKey, config.Port).Save(); err != nil {
		fmt.Println(err)
	}
	return
}
