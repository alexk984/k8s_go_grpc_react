package config

import (
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	DatabaseURL string
	GRPCPort    string
	HTTPPort    string
	JWTSecret   string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/k8s_grpc_db?sslmode=disable"),
		GRPCPort:    getEnv("GRPC_PORT", "8080"),
		HTTPPort:    getEnv("HTTP_PORT", "8081"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
