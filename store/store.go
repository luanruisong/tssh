package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const EnvName = "TSSH_HOME"

var configPath string

func DefaultCheck() error {
	configPath = os.Getenv(EnvName)
	if len(configPath) == 0 {
		return errors.New("env TSSH_HOME can not find")
	}

	if !fileExists(configPath) {
		return os.MkdirAll(configPath, os.ModePerm)
	}
	return nil
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

func ConfigExists(name string) bool {
	return fileExists(path.Join(configPath, name))
}

func Get(name string) (*SSHConfig, error) {
	finalPath := path.Join(configPath, name)
	if !fileExists(finalPath) {
		return nil, fmt.Errorf("config %s not exists", name)
	}
	return GetFromPath(finalPath)
}

func Del(name string) error {
	finalPath := path.Join(configPath, name)
	if !fileExists(finalPath) {
		return fmt.Errorf("config %s not exists", name)
	}
	err := os.Remove(finalPath)
	if err == nil {
		fmt.Println("delete", name, "success")
	}
	return err
}

func Set(cfg *SSHConfig) error {
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
