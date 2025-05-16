package config

import (
	"os"
	"strconv"
	"time"
)

// Config contains application configuration vars
type Config struct {
	ServerPort   int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Load loads configuration from env vars
func Load() *Config {
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	readTimeout, _ := strconv.Atoi(getEnv("READ_TIMEOUT", "15"))
	writeTimeout, _ := strconv.Atoi(getEnv("WRITE_TIMEOUT", "15"))
	idleTimeout, _ := strconv.Atoi(getEnv("IDLE_TIMEOUT", "60"))

	return &Config{
		ServerPort:   serverPort,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
