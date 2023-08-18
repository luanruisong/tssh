package store

import (
	"fmt"
	"os"
	"path"

	"github.com/luanruisong/tssh/constant"
)

func Env() string {

	return cfgPath
}

func ConfigExists(name string) bool {
	return fileExists(path.Join(Env(), name))
}

func Del(name string) error {
	finalPath := path.Join(Env(), name)
	if !fileExists(finalPath) {
		return fmt.Errorf("config %s not exists", name)
	}
	err := os.Remove(finalPath)
	if err == nil {
		fmt.Println("delete", name, "success")
	}
	return err
}

func FmtEnv() {
	fmt.Println(constant.EnvName, "=", Env())
}
