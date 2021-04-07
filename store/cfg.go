package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"text/tabwriter"
	"time"

	"tssh/ssh"

	ssh1 "golang.org/x/crypto/ssh"
)

type (
	SSHConfig struct {
		Name   string
		Ip     string
		User   string
		Pwd    string
		SshKey string
		Port   int
		SaveAt string
	}

	BatchConfig struct {
		list []*SSHConfig
		m    map[string]*SSHConfig
	}
)

func NewConfig(name, ip, user, pwd, sshKey string, port int) *SSHConfig {
	return &SSHConfig{
		Name:   name,
		Ip:     ip,
		User:   user,
		Pwd:    pwd,
		SshKey: sshKey,
		Port:   port,
		SaveAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (s *SSHConfig) SaveToPath(path string) error {
	b, e := json.MarshalIndent(s, "", "	")
	if e != nil {
		return e
	}
	err := ioutil.WriteFile(path, b, os.ModePerm)
	if err == nil {
		fmt.Println("save", s.Name, " success")
	}
	return err
}

func GetFromPath(path string) (s *SSHConfig, e error) {
	var b []byte
	b, e = ioutil.ReadFile(path)
	if e != nil {
		return
	}
	s = &SSHConfig{}
	e = json.Unmarshal(b, s)
	return
}

func (s *SSHConfig) Conn() (err error) {
	var (
		cfg *ssh1.ClientConfig
		cli *ssh1.Client
	)
	if len(s.SshKey) > 0 {
		cfg, err = ssh.PkCfg(s.User, s.SshKey)
	} else {
		cfg = ssh.PwdCfg(s.User, s.Pwd)
	}
	cli, err = ssh.Connect(s.Ip, s.Port, cfg)
	if err != nil {
		return err
	}
	return ssh.RunTerminal(cli, os.Stdin, os.Stdout, os.Stderr)
}

func (s *SSHConfig) ConnMode() string {
	if len(s.SshKey) > 0 {
		return "private_key"
	}
	return "password"
}

func (bc *BatchConfig) Println() {
	w := tabwriter.NewWriter(os.Stdout, 10, 3, 5, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "No\tname\tip\tuser\tauth_type\tport\tsave_at\t")
	for i, v := range bc.list {
		s := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%d\t%s\t", i+1, v.Name, v.Ip, v.User, v.ConnMode(), v.Port, v.SaveAt)
		fmt.Fprintln(w, s)
	}
	w.Flush()
}

func (bc *BatchConfig) Load() error {
	if err := DefaultCheck(); err != nil {
		return err
	}
	dir, err := ioutil.ReadDir(configPath)
	if err != nil {
		return err
	}
	list := make([]*SSHConfig, 0)
	m := make(map[string]*SSHConfig)
	for _, v := range dir {
		cfg := &SSHConfig{}
		b, e := ioutil.ReadFile(path.Join(configPath, v.Name()))
		if e != nil {
			return e
		}
		if e = json.Unmarshal(b, cfg); err == nil {
			list = append(list, cfg)
			m[cfg.Name] = cfg
		}
	}
	bc.list = list
	bc.m = m
	return nil
}

func (bc *BatchConfig) Get(str string) *SSHConfig {
	x, err := strconv.Atoi(str)
	if err != nil {
		return bc.m[str]
	}
	if x > len(bc.list) {
		return bc.m[str]
	}
	return bc.list[x-1]
}
