package util

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func GetIntEnv(key string, defaultValue int) int {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		log.WithFields(log.Fields{
			log.ErrorKey: err,
			"key":        key,
		}).Error("failed to get int env, use default")

		return defaultValue
	}

	return result
}

func GetStringEnv(key string, defaultValue string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	return value
}

func GetDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	result, err := time.ParseDuration(value)
	if err != nil {
		log.WithFields(log.Fields{
			log.ErrorKey: err,
			"key":        key,
		}).Error("failed to get duration env, use default")

		return defaultValue
	}

	return result
}
