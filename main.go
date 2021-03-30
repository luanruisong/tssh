package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"tssh/store"
)

func main() {

	if err := store.DefaultCheck(); err != nil {
		fmt.Println(err)
		return
	}

	var (
		a = flag.String("a", "", "add config {user@host}")
		s = flag.String("s", "", "set config {user@host}")
		d = flag.String("d", "", "del config {name}")
		c = flag.String("c", "", "connect config host {name}")
		l = flag.Bool("l", false, "config list")
		e = flag.Bool("e", false, "evn info")
	)

	var (
		n = flag.String("n", "", "set name in (-a|-s)")
		p = flag.String("p", "", "set password in (-a|-s)")
		P = flag.Int("P", 22, "set port in (-a|-s)")
		k = flag.String("k", "", "set private_key path in (-a|-s)")
	)

	flag.Parse()

	switch {
	case *l:
		list, err := store.List()
		if err != nil {
			fmt.Println(err)
			return
		}
		printCfg(list)
	case *e:
		store.Env()
	case len(*c) > 0:
		connByName(*c)
	case len(*d) > 0:
		err := store.Del(*d)
		if err != nil {
			fmt.Println(err)
			return
		}
	case len(*a) > 0:
		name := *n
		if len(name) == 0 {
			fmt.Println("config", name, "exists")
			_, name = GetUserAndHost(a)
		}
		if store.ConfigExists(name) {
			fmt.Println("config", name, "exists")
			return
		}
		fallthrough
	case len(*s) > 0:
		user, host := GetUserAndHost(a, s)
		if len(user) == 0 || len(host) == 0 {
			fmt.Println("user and host required")
			return
		}
		name := *n
		if len(name) == 0 {
			name = host
		}
		pwd, sshkey := *p, *k
		if len(pwd) == 0 && len(sshkey) == 0 {
			fmt.Println("pwd and sshkey required")
			return
		}
		cfg := store.NewConfig(name, host, user, pwd, sshkey, *P)
		err := store.Set(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("args not support")
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

func printCfg(cfgs []store.SSHConfig) {
	w := tabwriter.NewWriter(os.Stdout, 10, 3, 5, ' ',
		tabwriter.AlignRight)
	fmt.Fprintln(w, "No\tname\tip\tuser\tpwd\tkey_path\tport\tsave_at\t")
	for i, v := range cfgs {
		s := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\t%d\t%s\t", i+1, v.Name, v.Ip, v.User, v.Pwd, v.SshKey, v.Port, v.SaveAt)
		fmt.Fprintln(w, s)
	}
	w.Flush()
}

func connByName(name string) {
	info, err := store.Get(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := info.Conn(); err != nil {
		fmt.Println(err)
	}
}
