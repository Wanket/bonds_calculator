package util

type IGlobalConfig interface {
	MoexClientQueueSize() int

	HTTPPort() int
}

const (
	moexClientQueueSize = "CONFIG_MOEX_CLIENT_QUEUE_SIZE"
	httpPort            = "CONFIG_HTTP_PORT"
)

type GlobalConfig struct {
	moexClientQueueSize int

	httpPort int
}

//nolint:gomnd
func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		moexClientQueueSize: GetIntEnv(moexClientQueueSize, 10),
		httpPort:            GetIntEnv(httpPort, 8080),
	}
}

func (config *GlobalConfig) MoexClientQueueSize() int {
	return config.moexClientQueueSize
}

func (config *GlobalConfig) HTTPPort() int {
	return config.httpPort
}
