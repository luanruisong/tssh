package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"tssh/store"
)

type (
	CmdSSHConfig struct {
		Name           string
		Port           int
		Pwd            string
		PrivateKeyPath string
	}

	handler func(in string) string
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

func Interactive(def string, h handler) error {
	fmt.Printf(def)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			return err
		}
		inputStr := strings.TrimSpace(input)
		switch inputStr {
		case "exit", "q", "quit":
			return nil
		default:
			if x := h(inputStr); len(x) > 0 {
				fmt.Printf("%s", x)
			} else {
				return nil
			}
		}

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
	if len(body) == 0 {
		const tag = "please input {user@host}:"
		_ = Interactive(tag, func(in string) string {
			if len(in) > 0 {
				body = in
				return ""
			}
			return tag
		})
	}
	//获取添加/覆盖配置需要的参数
	config := ParseArgs(args)
	//检查别名输入情况
	if len(config.Name) == 0 {
		const tag = "please input alias name:"
		_ = Interactive(tag, func(in string) string {
			if len(in) > 0 {
				config.Name = in
			} else {
				config.Name = body
			}
			return ""
		})
	}
	if isAdd {
		if store.ConfigExists(config.Name) {
			fmt.Println("config", config.Name, "exists")
			return
		}
	}
	//检查密码和密钥输入情况
	if len(config.Pwd) == 0 && len(config.PrivateKeyPath) == 0 {
		const tag = "please input password:"
		_ = Interactive(tag, func(in string) string {
			if len(in) > 0 {
				config.Pwd = in
				return ""
			}
			return tag
		})
	}
	user, host := GetUserAndHost(body)
	if len(user) == 0 || len(host) == 0 {
		fmt.Println("user and host required")
		return
	}
	cfg := store.NewConfig(config.Name, host, user, config.Pwd, config.PrivateKeyPath, config.Port)
	err := store.Set(cfg)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func List() {
	batch := store.NewBatchConfig()
	if err := batch.Load(); err != nil {
		fmt.Println(err)
		return
	}
	const tag = "connect with index or name:"
	batch.Println()
	_ = Interactive(tag, func(in string) string {
		info := batch.Get(in)
		if info == nil {
			fmt.Println("can not find config", in)
			return tag
		}
		if err := info.Conn(); err != nil {
			fmt.Println(err)
		}
		return ""
	})
}

func Conn(name string) {
	batch := store.NewBatchConfig()
	if err := batch.Load(); err != nil {
		fmt.Println(err)
		return
	}
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
		fmt.Println("input alias name")
		return
	}
	err := store.Del(name)
	if err != nil {
		fmt.Println(err)
	}
}
