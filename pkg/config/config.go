package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
	// Добавь другие поля, если нужно (например SMTP, PORT и т.п.)
}

var (
	config *Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		// Загружаем .env если есть (игнорируем ошибку, если файла нет)
		_ = godotenv.Load()

		config = &Config{
			DBUrl:     mustEnv("DATABASE_URL"),
			JWTSecret: mustEnv("JWT_SECRET"),
			// тут можно добавить и другие конфиги
		}
	})
	return config
}

// Достаём переменную окружения или паникуем, если не задана
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return v
}
