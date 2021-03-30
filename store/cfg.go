package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"tssh/ssh"

	ssh1 "golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Name   string
	Ip     string
	User   string
	Pwd    string
	SshKey string
	Port   int
	SaveAt string
}

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
