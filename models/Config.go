package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func InitConfig() Config {

	err := godotenv.Load()
	if err != nil {
		if os.Getenv("DB_NAME") == "" {
			fmt.Println("Error loading conf! Create .env file or export os Environments")
		}
	}

	return Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}
