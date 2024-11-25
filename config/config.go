package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerPort string
	AppConfig  AppConfig
}

func LoadConfig() (*Config, error) {
	config := Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AppConfig:  getAppConfig(),
	}

	if config.AppConfig.DBHost == "" || config.AppConfig.DBPort == "" || config.AppConfig.DBUser == "" ||
		config.AppConfig.DBPassword == "" || config.AppConfig.DBName == "" {
		return nil, errors.New("database environment variables must be set")
	}

	if config.AppConfig.FromEmail == "" || config.AppConfig.FromEmailSMTP == "" || config.AppConfig.FromEmailPassword == "" ||
		config.AppConfig.SMTPAddrress == "" {
		return nil, errors.New("email environment variables must be set")
	}

	return &config, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getAppConfig() AppConfig {
	return AppConfig{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		FromEmail:         os.Getenv("FROM_EMAIL"),
		FromEmailPassword: os.Getenv("FROM_EMAIL_PASSWORD"),
		FromEmailSMTP:     os.Getenv("FROM_EMAIL_SMTP"),
		SMTPAddrress:      os.Getenv("SMTP_ADDR"),
	}
}
