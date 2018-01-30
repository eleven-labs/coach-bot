package config

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// Getenv is a proxy function from os.Getenv() but also logs in case of error.
func Getenv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Errorf("Unable to load environment variable '%s'", key)
		os.Exit(1)
	}

	return value
}
