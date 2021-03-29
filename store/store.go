package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const EnvName = "TSSH_HOME"

var configPath string

func init() {
	configPath = os.Getenv(EnvName)
	if len(configPath) == 0 {
		panic(errors.New("env TSSH_HOME can not find"))
	}

	if !fileExists(configPath) {
		os.MkdirAll(configPath, os.ModePerm)
	}
}

func fileExists(p string) bool {
	_, err := os.Stat(p) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Add(args []string) error {
	cfg, err := parseCfg(args)
	if err != nil {
		return err
	}
	finalPath := path.Join(configPath, cfg.Name)
	if fileExists(finalPath) {
		return fmt.Errorf("config %s exists", cfg.Name)
	}
	return cfg.SaveToPath(finalPath)
}

func GetByArgs(args []string) (*SSHConfig, error) {
	name, err := parseName(args)
	if err != nil {
		return nil, err
	}
	return Get(name)
}
func Get(name string) (*SSHConfig, error) {
	finalPath := path.Join(configPath, name)
	if !fileExists(finalPath) {
		return nil, fmt.Errorf("config %s not exists", name)
	}
	return GetFromPath(finalPath)
}

func Del(args []string) error {
	name, err := parseName(args)
	if err != nil {
		return err
	}
	finalPath := path.Join(configPath, name)
	if !fileExists(finalPath) {
		return fmt.Errorf("config %s not exists", name)
	}
	err = os.Remove(finalPath)
	if err == nil {
		fmt.Println("delete", name, "success")
	}
	return err
}

func Set(args []string) error {
	cfg, err := parseCfg(args)
	if err != nil {
		return err
	}
	finalPath := path.Join(configPath, cfg.Name)
	if fileExists(finalPath) {
		_ = os.Remove(finalPath)
	}
	return cfg.SaveToPath(finalPath)
}

func List() ([]SSHConfig, error) {
	dir, err := ioutil.ReadDir(configPath)
	if err != nil {
		return nil, err
	}
	res := make([]SSHConfig, 0)
	for _, v := range dir {
		cfg := SSHConfig{}
		b, e := ioutil.ReadFile(path.Join(configPath, v.Name()))
		if e != nil {
			return nil, e
		}
		if e = json.Unmarshal(b, &cfg); err == nil {
			res = append(res, cfg)
		} else {
			return nil, e
		}
	}
	return res, nil
}

func Env() {
	fmt.Println("env", EnvName, "=", configPath)
}

func parseCfg(args []string) (*SSHConfig, error) {
	if len(args) < 6 {
		return nil, errors.New("args error")
	}
	name, ip, user, pwd := args[2], args[3], args[4], args[5]
	port := 22
	if len(args) > 6 {
		port, _ = strconv.Atoi(args[6])
	}
	return NewConfig(name, ip, user, pwd, port), nil
}

func parseName(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("can not get name")
	}
	return args[2], nil
}
