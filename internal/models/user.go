package models

import (
	"time"

	"gorm.io/gorm"
)

// User представляет модель пользователя в базе данных
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Name         string         `gorm:"not null;size:255" json:"name"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string         `gorm:"not null;size:255" json:"-"` // Хеш пароля, не возвращается в JSON
	Role         string         `gorm:"not null;default:'user';size:50" json:"role"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName возвращает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}
