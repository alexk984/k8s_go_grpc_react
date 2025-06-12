package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKey тип для ключей контекста
type ContextKey string

const (
	// UserContextKey ключ для пользователя в контексте
	UserContextKey ContextKey = "user"
)

// AuthMiddleware middleware для проверки JWT токенов
type AuthMiddleware struct {
	jwtService JWTService
	logger     *logrus.Logger
}

// NewAuthMiddleware создает новый экземпляр middleware
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: NewJWTService(),
		logger:     logrus.New(),
	}
}

// NewAuthMiddlewareWithDeps создает новый экземпляр AuthMiddleware с зависимостями
func NewAuthMiddlewareWithDeps(jwtService JWTService, logger *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// UnaryInterceptor возвращает unary interceptor для аутентификации
func (m *AuthMiddleware) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Публичные методы, не требующие аутентификации
	publicMethods := map[string]bool{
		"/user.UserService/Register": true,
		"/user.UserService/Login":    true,
	}

	// Если метод публичный, пропускаем аутентификацию
	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// Извлекаем токен из метаданных
	token, err := m.extractTokenFromMetadata(ctx)
	if err != nil {
		m.logger.WithError(err).Warn("Ошибка извлечения токена")
		return nil, status.Error(codes.Unauthenticated, "Токен не предоставлен")
	}

	// Валидируем токен
	claims, err := m.jwtService.ValidateToken(token)
	if err != nil {
		m.logger.WithError(err).Warn("Недействительный токен")
		return nil, status.Error(codes.Unauthenticated, "Недействительный токен")
	}

	// Добавляем информацию о пользователе в контекст
	ctx = context.WithValue(ctx, "user_id", claims.UserID)
	ctx = context.WithValue(ctx, "user_email", claims.Email)
	ctx = context.WithValue(ctx, "user_role", claims.Role)

	return handler(ctx, req)
}

// extractTokenFromMetadata извлекает JWT токен из gRPC метаданных
func (m *AuthMiddleware) extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Метаданные не найдены")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", status.Error(codes.Unauthenticated, "Заголовок Authorization не найден")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "Неверный формат токена")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// RequireAuth middleware для обязательной аутентификации
func RequireAuth(jwtService JWTService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Извлекаем токен из метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Метаданные не найдены")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Токен не предоставлен")
		}

		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "Неверный формат токена")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Недействительный токен")
		}

		// Добавляем информацию о пользователе в контекст
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		return handler(ctx, req)
	}
}

// OptionalAuth middleware для опциональной аутентификации
func OptionalAuth(jwtService JWTService) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Пытаемся извлечь токен из метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req) // Продолжаем без аутентификации
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return handler(ctx, req) // Продолжаем без аутентификации
		}

		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return handler(ctx, req) // Продолжаем без аутентификации
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return handler(ctx, req) // Продолжаем без аутентификации при ошибке
		}

		// Добавляем информацию о пользователе в контекст
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		return handler(ctx, req)
	}
}

// RequireRole middleware для проверки роли пользователя
func RequireRole(jwtService JWTService, requiredRole string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Сначала проверяем аутентификацию
		authInterceptor := RequireAuth(jwtService)
		_, err := authInterceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil // Заглушка
		})
		if err != nil {
			return nil, err
		}

		// Проверяем роль
		userRole, ok := ctx.Value("user_role").(string)
		if !ok || userRole != requiredRole {
			return nil, status.Error(codes.PermissionDenied, "Недостаточно прав доступа")
		}

		return handler(ctx, req)
	}
}

// RequireAuth middleware для обязательной аутентификации
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			m.logger.WithField("path", r.URL.Path).Warn("Missing authorization token")
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.WithError(err).WithField("path", r.URL.Path).Warn("Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware для опциональной аутентификации
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token != "" {
			claims, err := m.jwtService.ValidateToken(token)
			if err == nil {
				// Если токен валидный, добавляем пользователя в контекст
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRole middleware для проверки роли пользователя
func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*Claims)
			if !ok {
				m.logger.WithField("path", r.URL.Path).Warn("User not authenticated")
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Проверяем роль пользователя
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				m.logger.WithFields(logrus.Fields{
					"user_id":        claims.UserID,
					"user_role":      claims.Role,
					"required_roles": roles,
					"path":           r.URL.Path,
				}).Warn("Insufficient permissions")
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken извлекает JWT токен из заголовка Authorization
func (m *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Ожидаем формат: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// GetUserFromContext извлекает информацию о пользователе из контекста
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*Claims)
	return claims, ok
}
