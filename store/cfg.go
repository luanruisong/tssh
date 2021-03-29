package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type SSHConfig struct {
	Name   string
	Ip     string
	User   string
	Pwd    string
	Port   int
	SaveAt string
}

func NewConfig(name, ip, user, pwd string, port int) *SSHConfig {
	return &SSHConfig{
		Name:   name,
		Ip:     ip,
		User:   user,
		Pwd:    pwd,
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
