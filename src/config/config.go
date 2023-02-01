package config

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Config struct {
	Debug bool
}

// LoadConfig loads configuration values from environment variables
func LoadConfig() *Config {
	config := Config{
		Debug: parseBool("DEBUG", true),
	}

	setLogging(config.Debug)
	return &config
}

func setLogging(debug bool) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debugging mode enabled")
	}
}

// parseBool parses an environment variable as a boolean value
func parseBool(key string, fallback bool) bool {
	if value := os.Getenv(key); strings.ToLower(value) == "true" {
		return true
	}
	if value := os.Getenv(key); strings.ToLower(value) == "false" {
		return false
	}
	return fallback
}
