package util

import (
	"sync"
)

var config GlobalConfig
var loadConfigOnce sync.Once

type GlobalConfig struct {
	moexClientQueueSize int

	httpPort int
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

func (config *GlobalConfig) HttpPort() int {
	return config.httpPort
}

const (
	moexClientQueueSize = "CONFIG_MOEX_CLIENT_QUEUE_SIZE"
	httpPort            = "CONFIG_HTTP_PORT"
)

func loadGlobalConfig() GlobalConfig {
	return GlobalConfig{
		moexClientQueueSize: GetIntEnv(moexClientQueueSize, 10),
		httpPort:            GetIntEnv(httpPort, 8080),
	}
}
