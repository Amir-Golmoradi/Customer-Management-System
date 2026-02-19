package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTP        HTTPConfig
	DB          DBConfig
	LogLevel    string
}

type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

type DBConfig struct {
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	PingTimeout     time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	cfg := &Config{
		DatabaseURL: databaseURL,
		HTTP: HTTPConfig{
			Port:            getEnv("HTTP_PORT", "8080"),
			ReadTimeout:     getDurationEnv("HTTP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout:    getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:     getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDurationEnv("HTTP_SHUTDOWN_TIMEOUT", 15*time.Second),
			MaxHeaderBytes:  getIntEnv("HTTP_MAX_HEADER_BYTES", 1<<20),
		},
		DB: DBConfig{
			MaxConns:        int32(getIntEnv("DB_MAX_CONNS", 20)),
			MinConns:        int32(getIntEnv("DB_MIN_CONNS", 2)),
			MaxConnLifetime: getDurationEnv("DB_MAX_CONN_LIFETIME", 30*time.Minute),
			MaxConnIdleTime: getDurationEnv("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),
			PingTimeout:     getDurationEnv("DB_PING_TIMEOUT", 2*time.Second),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
