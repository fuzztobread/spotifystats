package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	RedisAddr    string
	KafkaBrokers []string
	Port         string
}

func Load() *Config {
	godotenv.Load(".env.local")
	godotenv.Load(".env")

	return &Config{
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		KafkaBrokers: []string{os.Getenv("KAFKA_BROKERS")},
		Port:         os.Getenv("PORT"),
	}
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" + c.DBHost + ":" + c.DBPort + "/" + c.DBName
}
