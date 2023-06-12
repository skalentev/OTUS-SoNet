package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetDBConfig(pref string) DBConfig {

	err := godotenv.Load()
	if err != nil {
		if os.Getenv("DB_NAME") == "" {
			fmt.Println("Error loading conf! Create .env file or export os Environments")
		}
	}

	return DBConfig{
		Driver:   os.Getenv(pref + "DRIVER"),
		Host:     os.Getenv(pref + "HOST"),
		Port:     os.Getenv(pref + "PORT"),
		User:     os.Getenv(pref + "USER"),
		Password: os.Getenv(pref + "PASSWORD"),
		DBName:   os.Getenv(pref + "NAME"),
		SSLMode:  os.Getenv(pref + "SSLMODE"),
	}
}
