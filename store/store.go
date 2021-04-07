package store

import (
	"fmt"
	"os"
	"path"
)

const EnvName = "TSSH_HOME"

var configPath string

func DefaultCheck() error {
	configPath = os.Getenv(EnvName)
	if len(configPath) == 0 {
		home := os.Getenv("HOME")
		if len(home) == 0 {
			return fmt.Errorf("env '%s' not found,please set a dir in env", EnvName)
		}
		configPath = path.Join(home, ".tssh/config")
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

func Del(name string) error {
	if err := DefaultCheck(); err != nil {
		return err
	}
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
	if err := DefaultCheck(); err != nil {
		return err
	}
	finalPath := path.Join(configPath, cfg.Name)
	if fileExists(finalPath) {
		_ = os.Remove(finalPath)
	}
	return cfg.SaveToPath(finalPath)
}

func Env() {
	_ = DefaultCheck()
	fmt.Println(EnvName, "=", configPath)
}

func NewBatchConfig() *BatchConfig {
	return &BatchConfig{}
}
