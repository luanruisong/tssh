package store

func GetBatchConfig() *configBatch {
	cacheOnce.Do(func() {
		global = &configBatch{}
		global.Load()
	})
	return global
}

func ListConfig() []*SSHConfig {
	return GetBatchConfig().List()
}

func GetConfig(name string) *SSHConfig {
	return GetBatchConfig().Get(name)
}
