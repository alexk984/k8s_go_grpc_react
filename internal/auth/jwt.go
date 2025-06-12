package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims представляет JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTService интерфейс для работы с JWT токенами
type JWTService interface {
	GenerateToken(userID uint, email, role string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
}

// jwtService реализация JWT сервиса
type jwtService struct {
	secretKey       []byte
	tokenExpiration time.Duration
}

// NewJWTService создает новый экземпляр JWT сервиса
func NewJWTService() JWTService {
	secretKey := getJWTSecret()
	expiration := getTokenExpiration()

	return &jwtService{
		secretKey:       []byte(secretKey),
		tokenExpiration: expiration,
	}
}

// GenerateToken создает новый JWT токен
func (j *jwtService) GenerateToken(userID uint, email, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "k8s-grpc-app",
			Subject:   fmt.Sprintf("user:%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет и парсит JWT токен
func (j *jwtService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// HashPassword хеширует пароль с использованием bcrypt
func (j *jwtService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// CheckPassword проверяет соответствие пароля хешу
func (j *jwtService) CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password mismatch: %w", err)
	}
	return nil
}

// getJWTSecret получает секретный ключ из переменных окружения
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// В production должен быть установлен через переменные окружения
		secret = "your-super-secret-jwt-key-change-in-production"
	}
	return secret
}

// getTokenExpiration получает время жизни токена из переменных окружения
func getTokenExpiration() time.Duration {
	expStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if expStr == "" {
		return 24 * time.Hour // По умолчанию 24 часа
	}

	hours, err := strconv.Atoi(expStr)
	if err != nil {
		return 24 * time.Hour
	}

	return time.Duration(hours) * time.Hour
}
