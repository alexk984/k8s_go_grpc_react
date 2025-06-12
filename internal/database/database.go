package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"k8s-go-grpc-react/internal/models"
)

// Config содержит конфигурацию базы данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetConfigFromEnv загружает конфигурацию базы данных из переменных окружения
func GetConfigFromEnv() *Config {
	port := 5432
	if portStr := getEnv("DB_PORT", "5432"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "k8s_grpc_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect подключается к базе данных по URL
func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// NewPostgresDB создает новое подключение к PostgreSQL
func NewPostgresDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	return Connect(dsn)
}

// Migrate выполняет миграции базы данных
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

// AutoMigrate выполняет автоматические миграции (для обратной совместимости)
func AutoMigrate(db *gorm.DB) error {
	return Migrate(db)
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
