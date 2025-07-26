package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName string
	Port    string
}

type DbConfig struct {
	Dsn string
}

type Config struct {
	App *AppConfig
	Db  *DbConfig
}

func NewConfig(envPath string) *Config {
	err := godotenv.Load(envPath)
	if err != nil {
		panic("Error loading .env file")
	}

	appName := os.Getenv("APP_NAME")
	port := ":" + os.Getenv("PORT")
	dsn := os.Getenv("DB_DSN")

	return &Config{
		App: &AppConfig{
			AppName: appName,
			Port:    port,
		},
		Db: &DbConfig{
			Dsn: dsn,
		},
	}
}
