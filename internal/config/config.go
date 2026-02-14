package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
}

// Since i don't want to read the memory address of each field
// I use Pointer of Config struct

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBName:      os.Getenv("DB_NAME"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
	}, nil
}
