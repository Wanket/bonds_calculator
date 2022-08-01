package util

import "time"

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type IGlobalConfig interface {
	MoexClientQueueSize() int

	HTTPPort() int

	GetDataBaseConfig() DBConfig
	GetRedisConfig() RedisConfig

	GetDefaultRequestTimeout() time.Duration

	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
}

const (
	moexClientQueueSize = "CONFIG_MOEX_CLIENT_QUEUE_SIZE"

	httpPort              = "CONFIG_HTTP_PORT"
	defaultRequestTimeout = "CONFIG_REQUEST_TIMEOUT"

	accessTokenTTL  = "CONFIG_ACCESS_TOKEN_TTL"
	refreshTokenTTL = "CONFIG_REFRESH_TOKEN_TTL" //nolint:gosec

	dbUser     = "CONFIG_DB_USER"
	dbPassword = "CONFIG_DB_PASSWORD" //nolint:gosec
	dbHost     = "CONFIG_DB_HOST"
	dbPort     = "CONFIG_DB_PORT"
	dbName     = "CONFIG_DB_NAME"

	redisHost     = "CONFIG_REDIS_HOST"
	redisPort     = "CONFIG_REDIS_PORT"
	redisPassword = "CONFIG_REDIS_PASSWORD" //nolint:gosec
	redisDB       = "CONFIG_REDIS_DB"
)

type GlobalConfig struct {
	moexClientQueueSize int

	httpPort int

	defaultRequestTimeout time.Duration

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	dbConfig    DBConfig
	redisConfig RedisConfig
}

//nolint:gomnd
func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		moexClientQueueSize:   GetIntEnv(moexClientQueueSize, 10),
		httpPort:              GetIntEnv(httpPort, 8080),
		defaultRequestTimeout: GetDurationEnv(defaultRequestTimeout, time.Millisecond*100),
		accessTokenTTL:        GetDurationEnv(accessTokenTTL, time.Minute*15),
		refreshTokenTTL:       GetDurationEnv(refreshTokenTTL, Week),
		dbConfig: DBConfig{
			Host:     GetStringEnv(dbHost, "localhost"),
			Port:     GetIntEnv(dbPort, 5432),
			User:     GetStringEnv(dbUser, "bonds_calculator"),
			Password: GetStringEnv(dbPassword, "bonds_calculator"),
			DBName:   GetStringEnv(dbName, "bonds_calculator"),
		},
		redisConfig: RedisConfig{
			Host:     GetStringEnv(redisHost, "localhost"),
			Port:     GetIntEnv(redisPort, 6379),
			Password: GetStringEnv(redisPassword, ""),
			DB:       GetIntEnv(redisDB, 0),
		},
	}
}

func (config *GlobalConfig) MoexClientQueueSize() int {
	return config.moexClientQueueSize
}

func (config *GlobalConfig) HTTPPort() int {
	return config.httpPort
}

func (config *GlobalConfig) GetDataBaseConfig() DBConfig {
	return config.dbConfig
}

func (config *GlobalConfig) GetRedisConfig() RedisConfig {
	return config.redisConfig
}

func (config *GlobalConfig) GetDefaultRequestTimeout() time.Duration {
	return config.defaultRequestTimeout
}

func (config *GlobalConfig) GetAccessTokenTTL() time.Duration {
	return config.accessTokenTTL
}

func (config *GlobalConfig) GetRefreshTokenTTL() time.Duration {
	return config.refreshTokenTTL
}
