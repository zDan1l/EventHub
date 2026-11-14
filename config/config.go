package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseURL string

	JWTSecret      string
	JWTExpireHours string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("config: file .env tidak ditemukan, memakai environment sistem")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "eventhub"),
		DBPassword:  getEnv("DB_PASSWORD", "eventhub"),
		DBName:      getEnv("DB_NAME", "eventhub"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpireHours: getEnv("JWT_EXPIRE_HOURS", "24"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
