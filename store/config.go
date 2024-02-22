package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/luanruisong/tssh/ssh"

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

func (s *SSHConfig) saveToPath(path string) error {
	b, e := json.MarshalIndent(s, "", "	")
	if e != nil {
		return e
	}
	err := os.WriteFile(path, b, os.ModePerm)
	if err == nil {
		fmt.Println(fmt.Sprintf("save <%s> success", s.Name))
	}
	return err
}

func (s *SSHConfig) Save() error {
	finalPath := path.Join(Env(), s.Name)
	if fileExists(finalPath) {
		_ = os.Remove(finalPath)
	}
	return s.saveToPath(finalPath)
}

func (s *SSHConfig) String() string {
	return s.Name
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
	if err != nil {
		return err
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
