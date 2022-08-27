package enviroment

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetEnvString(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Printf("empty env: %s, default: %s", envName, defaultValue)
		return defaultValue
	}
	return value
}

func GetEnvDuration(envName string, defaultValue time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(envName))
	if err != nil {
		log.Printf("error of env %s: %s, default: %v", envName, err.Error(), defaultValue)
		return defaultValue
	}
	return value
}

func GetEnvBool(envName string, defaultValue bool) bool {
	value, err := strconv.ParseBool(os.Getenv(envName))
	if err != nil {
		log.Printf("error of env %s: %s, default: %v", envName, err.Error(), defaultValue)
		return defaultValue
	}
	return value
}
