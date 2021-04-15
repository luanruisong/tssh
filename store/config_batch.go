package store

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type (
	configBatch struct {
		list []*SSHConfig
		m    map[string]*SSHConfig
	}
)

func (bc *configBatch) Load() {
	if len(bc.list) > 0 && len(bc.m) > 0 {
		return
	}
	env := Env()
	dir, err := ioutil.ReadDir(env)
	if err != nil {
		return
	}
	list := make([]*SSHConfig, 0)
	m := make(map[string]*SSHConfig)
	for _, v := range dir {
		cfg := &SSHConfig{}
		var b []byte
		if b, err = ioutil.ReadFile(path.Join(env, v.Name())); err != nil {
			return
		}
		if err = json.Unmarshal(b, cfg); err == nil {
			cfg.Name = v.Name()
			list = append(list, cfg)
			m[cfg.Name] = cfg
		}
	}
	bc.list = list
	bc.m = m
}

func (bc *configBatch) Get(str string) *SSHConfig {
	return bc.m[str]
}

func (bc *configBatch) List() []*SSHConfig {
	return bc.list
}
