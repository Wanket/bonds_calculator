package util

import (
	"sync"
)

var config GlobalConfig
var loadConfigOnce sync.Once

type GlobalConfig struct {
	moexClientQueueSize int
}

func GetGlobalConfig() *GlobalConfig {
	loadConfigOnce.Do(func() {
		config = loadGlobalConfig()
	})

	return &config
}

func (config *GlobalConfig) MoexClientQueueSize() int {
	return config.moexClientQueueSize
}

const (
	moexClientQueueSize = "CONFIG_MOEX_CLIENT_QUEUE_SIZE"
)

func loadGlobalConfig() GlobalConfig {
	return GlobalConfig{
		moexClientQueueSize: GetIntEnv(moexClientQueueSize, 10),
	}
}
