package config

import (
	"os"
)

type Config struct {
	DBURL      string
	RabbitMQURL string
	RedisURL   string
}

func LoadConfig() *Config {
	return &Config{
		DBURL:       os.Getenv("DB_URL"),
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
	}
}
