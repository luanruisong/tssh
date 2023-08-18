package store

import (
	"os"
	"path"
	"sync"

	"github.com/luanruisong/tssh/constant"
)

var (
	cfgPath   string
	global    *configBatch
	cacheOnce *sync.Once
)

func init() {
	cacheOnce = &sync.Once{}
	cfgPath = os.Getenv(constant.EnvName)
	if len(cfgPath) == 0 {
		cfgPath = buildConfigPath()
	}
	if !fileExists(cfgPath) {
		_ = os.MkdirAll(cfgPath, os.ModePerm)
	}
}

func buildConfigPath() string {
	return path.Join(os.Getenv(constant.HOME), ".tssh", "config")
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
