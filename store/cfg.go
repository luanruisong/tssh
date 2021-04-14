package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"tssh/ssh"

	ssh1 "golang.org/x/crypto/ssh"
)

type (
	SSHConfig struct {
		Name   string `json:"-"`
		Ip     string
		User   string
		Pwd    string
		SshKey []byte
		Port   int
		SaveAt string
	}

	BatchConfig struct {
		list []*SSHConfig
		m    map[string]*SSHConfig
	}
)

func NewConfig(name, ip, user, pwd string, sshKey []byte, port int) *SSHConfig {
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

func (s *SSHConfig) FmtName() string {
	name := s.Name
	if len(name) > 20 {
		name = name[:17] + "..."
	}
	return fmt.Sprintf("%-20s", name)
}

func (bc *BatchConfig) Load() error {
	if len(bc.list) > 0 && len(bc.m) > 0 {
		return nil
	}
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
		var b []byte
		if b, err = ioutil.ReadFile(path.Join(configPath, v.Name())); err != nil {
			return err
		}
		if err = json.Unmarshal(b, cfg); err == nil {
			cfg.Name = v.Name()
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

func (bc *BatchConfig) List() []*SSHConfig {
	return bc.list
}
