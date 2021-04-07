package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"tssh/store"
)

var (
	n = flag.String("n", "", "set name in (-a|-s)")
	p = flag.String("p", "", "set password in (-a|-s)")
	P = flag.Int("P", 22, "set port in (-a|-s)")
	k = flag.String("k", "", "set private_key path in (-a|-s)")
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

func ParseName(n *string, a ...*string) string {
	if len(*n) > 0 {
		return *n
	}
	for i := range a {
		if len(*a[i]) > 0 {
			return *a[i]
		}
	}
	return ""
}

func ParseArgs() *CmdSSHConfig {

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

func GetUserAndHost(a ...*string) (string, string) {
	for i := range a {
		curr := *a[i]
		if len(curr) > 0 {
			if idx := strings.Index(curr, "@"); idx > 0 {
				return curr[:idx], curr[idx+1:]
			}
		}
	}
	return "", ""
}

func AddOrSave(a, s *string) (f bool) {
	if len(*a) == 0 && len(*s) == 0 {
		return
	}
	f = true
	//获取添加/覆盖配置需要的参数
	config := ParseArgs()
	//检查别名输入情况
	if len(config.Name) == 0 {
		const tag = "please input alias name:"
		_ = Interactive(tag, func(in string) string {
			if len(in) > 0 {
				config.Name = in
			} else {
				config.Name = ParseName(a, s)
			}
			return ""
		})
	}
	if len(*a) > 0 {
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
	user, host := GetUserAndHost(a, s)
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

func ListAndConn(l *bool, c *string) (f bool) {
	if !*l && len(*c) == 0 {
		return
	}
	f = true
	batch := store.NewBatchConfig()
	if err := batch.Load(); err != nil {
		fmt.Println(err)
		return
	}
	if *l {
		const tag = "conn with index or name:"
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

	if len(*c) > 0 {
		name := ParseName(c)
		info := batch.Get(name)
		if info == nil {
			fmt.Println("can not find config", name)
			return
		}
		if err := info.Conn(); err != nil {
			fmt.Println(err)
		}
	}
	return
}
