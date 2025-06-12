package service

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"k8s-go-grpc-react/internal/auth"
	"k8s-go-grpc-react/internal/models"
	"k8s-go-grpc-react/internal/repository"
	pb "k8s-go-grpc-react/proto"
)

// UserService реализует gRPC сервис для работы с пользователями
type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo        repository.UserRepository
	jwtService      auth.JWTService
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	usersCount      prometheus.Gauge
}

// NewUserService создает новый экземпляр UserService без метрик
func NewUserService(userRepo repository.UserRepository) *UserService {
	return NewUserServiceWithMetrics(userRepo, nil, nil, nil)
}

// NewUserServiceWithMetrics создает новый экземпляр UserService с метриками
func NewUserServiceWithMetrics(userRepo repository.UserRepository, requestsTotal *prometheus.CounterVec, requestDuration *prometheus.HistogramVec, usersCount prometheus.Gauge) *UserService {
	service := &UserService{
		userRepo:        userRepo,
		jwtService:      auth.NewJWTService(),
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
		usersCount:      usersCount,
	}

	// Обновляем счетчик пользователей
	if usersCount != nil {
		go func() {
			if count, err := userRepo.Count(context.Background()); err == nil {
				usersCount.Set(float64(count))
			}
		}()
	}

	return service
}

// recordMetrics записывает метрики для запроса
func (s *UserService) recordMetrics(method string, status string, duration time.Duration) {
	if s.requestsTotal != nil {
		s.requestsTotal.WithLabelValues(method, status).Inc()
	}
	if s.requestDuration != nil {
		s.requestDuration.WithLabelValues(method).Observe(duration.Seconds())
	}
}

// modelToProto конвертирует модель пользователя в protobuf
func (s *UserService) modelToProto(user *models.User) *pb.User {
	return &pb.User{
		Id:        int32(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Unix(),
	}
}

// Register регистрирует нового пользователя
func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	start := time.Now()
	defer func() {
		s.recordMetrics("Register", "success", time.Since(start))
	}()

	if req.Name == "" {
		s.recordMetrics("Register", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Имя пользователя не может быть пустым")
	}

	if req.Email == "" {
		s.recordMetrics("Register", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Email не может быть пустым")
	}

	if req.Password == "" {
		s.recordMetrics("Register", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Пароль не может быть пустым")
	}

	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.recordMetrics("Register", "already_exists", time.Since(start))
		return nil, status.Error(codes.AlreadyExists, "Пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := s.jwtService.HashPassword(req.Password)
	if err != nil {
		s.recordMetrics("Register", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при хешировании пароля")
	}

	// Создаем нового пользователя
	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user", // По умолчанию роль пользователя
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.recordMetrics("Register", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при создании пользователя")
	}

	// Генерируем JWT токен
	token, err := s.jwtService.GenerateToken(newUser.ID, newUser.Email, newUser.Role)
	if err != nil {
		s.recordMetrics("Register", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при генерации токена")
	}

	// Обновляем счетчик пользователей
	if s.usersCount != nil {
		go func() {
			if count, err := s.userRepo.Count(context.Background()); err == nil {
				s.usersCount.Set(float64(count))
			}
		}()
	}

	return &pb.AuthResponse{
		Token:   token,
		User:    s.modelToProto(newUser),
		Message: "Пользователь успешно зарегистрирован",
	}, nil
}

// Login выполняет вход пользователя в систему
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	start := time.Now()
	defer func() {
		s.recordMetrics("Login", "success", time.Since(start))
	}()

	if req.Email == "" {
		s.recordMetrics("Login", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Email не может быть пустым")
	}

	if req.Password == "" {
		s.recordMetrics("Login", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Пароль не может быть пустым")
	}

	// Ищем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.recordMetrics("Login", "not_found", time.Since(start))
		return nil, status.Error(codes.NotFound, "Неверный email или пароль")
	}

	// Проверяем активность пользователя
	if !user.IsActive {
		s.recordMetrics("Login", "permission_denied", time.Since(start))
		return nil, status.Error(codes.PermissionDenied, "Аккаунт заблокирован")
	}

	// Проверяем пароль
	if err := s.jwtService.CheckPassword(user.PasswordHash, req.Password); err != nil {
		s.recordMetrics("Login", "unauthenticated", time.Since(start))
		return nil, status.Error(codes.Unauthenticated, "Неверный email или пароль")
	}

	// Генерируем JWT токен
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		s.recordMetrics("Login", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при генерации токена")
	}

	return &pb.AuthResponse{
		Token:   token,
		User:    s.modelToProto(user),
		Message: "Успешный вход в систему",
	}, nil
}

// GetUser получает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	start := time.Now()
	defer func() {
		s.recordMetrics("GetUser", "success", time.Since(start))
	}()

	user, err := s.userRepo.GetByID(ctx, uint(req.Id))
	if err != nil {
		s.recordMetrics("GetUser", "not_found", time.Since(start))
		return nil, status.Error(codes.NotFound, "Пользователь не найден")
	}

	return &pb.UserResponse{
		User:    s.modelToProto(user),
		Message: "Пользователь успешно найден",
	}, nil
}

// CreateUser создает нового пользователя (только для админов)
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	start := time.Now()
	defer func() {
		s.recordMetrics("CreateUser", "success", time.Since(start))
	}()

	if req.Name == "" {
		s.recordMetrics("CreateUser", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Имя пользователя не может быть пустым")
	}

	if req.Email == "" {
		s.recordMetrics("CreateUser", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Email не может быть пустым")
	}

	if req.Password == "" {
		s.recordMetrics("CreateUser", "invalid_argument", time.Since(start))
		return nil, status.Error(codes.InvalidArgument, "Пароль не может быть пустым")
	}

	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.recordMetrics("CreateUser", "already_exists", time.Since(start))
		return nil, status.Error(codes.AlreadyExists, "Пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := s.jwtService.HashPassword(req.Password)
	if err != nil {
		s.recordMetrics("CreateUser", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при хешировании пароля")
	}

	// Устанавливаем роль
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Создаем нового пользователя
	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.recordMetrics("CreateUser", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при создании пользователя")
	}

	// Обновляем счетчик пользователей
	if s.usersCount != nil {
		go func() {
			if count, err := s.userRepo.Count(context.Background()); err == nil {
				s.usersCount.Set(float64(count))
			}
		}()
	}

	return &pb.UserResponse{
		User:    s.modelToProto(newUser),
		Message: "Пользователь успешно создан",
	}, nil
}

// ListUsers возвращает список всех пользователей
func (s *UserService) ListUsers(ctx context.Context, req *pb.Empty) (*pb.UserListResponse, error) {
	start := time.Now()
	defer func() {
		s.recordMetrics("ListUsers", "success", time.Since(start))
	}()

	users, err := s.userRepo.List(ctx, 0, 0) // Без ограничений
	if err != nil {
		s.recordMetrics("ListUsers", "internal_error", time.Since(start))
		return nil, status.Error(codes.Internal, "Ошибка при получении списка пользователей")
	}

	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = s.modelToProto(user)
	}

	return &pb.UserListResponse{
		Users: protoUsers,
		Total: int32(len(protoUsers)),
	}, nil
}
