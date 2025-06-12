package service

import (
	"context"
	"testing"

	"k8s-go-grpc-react/internal/models"
	pb "k8s-go-grpc-react/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository - мок репозитория для тестирования
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}

	// Настраиваем мок - пользователь с таким email не существует
	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, assert.AnError)

	// Настраиваем мок для создания пользователя
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	// Act
	resp, err := service.CreateUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Role, resp.User.Role)
	assert.True(t, resp.User.IsActive)
	assert.Equal(t, "Пользователь успешно создан", resp.Message)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(1)
	expectedUser := &models.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Role:     "user",
		IsActive: true,
	}

	// Настраиваем мок для получения пользователя
	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	req := &pb.GetUserRequest{Id: int32(userID)}

	// Act
	resp, err := service.GetUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser.Name, resp.User.Name)
	assert.Equal(t, expectedUser.Email, resp.User.Email)
	assert.Equal(t, expectedUser.Role, resp.User.Role)
	assert.True(t, resp.User.IsActive)
	assert.Equal(t, "Пользователь успешно найден", resp.Message)

	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	expectedUsers := []*models.User{
		{
			ID:       1,
			Name:     "User 1",
			Email:    "user1@example.com",
			Role:     "user",
			IsActive: true,
		},
		{
			ID:       2,
			Name:     "User 2",
			Email:    "user2@example.com",
			Role:     "admin",
			IsActive: true,
		},
	}

	// Настраиваем мок для получения списка пользователей
	mockRepo.On("List", ctx, 0, 0).Return(expectedUsers, nil)

	req := &pb.Empty{}

	// Act
	resp, err := service.ListUsers(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, len(expectedUsers))
	assert.Equal(t, int32(len(expectedUsers)), resp.Total)

	for i, user := range resp.Users {
		assert.Equal(t, expectedUsers[i].Name, user.Name)
		assert.Equal(t, expectedUsers[i].Email, user.Email)
		assert.Equal(t, expectedUsers[i].Role, user.Role)
		assert.True(t, user.IsActive)
	}

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_InvalidInput(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()

	testCases := []struct {
		name        string
		req         *pb.CreateUserRequest
		expectedErr string
	}{
		{
			name: "Empty name",
			req: &pb.CreateUserRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedErr: "Имя пользователя не может быть пустым",
		},
		{
			name: "Empty email",
			req: &pb.CreateUserRequest{
				Name:     "Test User",
				Email:    "",
				Password: "password123",
			},
			expectedErr: "Email не может быть пустым",
		},
		{
			name: "Empty password",
			req: &pb.CreateUserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			expectedErr: "Пароль не может быть пустым",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			resp, err := service.CreateUser(ctx, tc.req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	email := "test@example.com"

	// Второй запрос на создание пользователя с существующим email
	req := &pb.CreateUserRequest{
		Name:     "User 2",
		Email:    email,
		Password: "password456",
		Role:     "admin",
	}

	// Настраиваем мок - пользователь с таким email уже существует
	existingUser := &models.User{
		ID:    1,
		Name:  "User 1",
		Email: email,
		Role:  "user",
	}
	mockRepo.On("GetByEmail", ctx, email).Return(existingUser, nil)

	// Act
	resp, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Пользователь с таким email уже существует")

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	ctx := context.Background()
	userID := uint(999)

	// Настраиваем мок - пользователь не найден
	mockRepo.On("GetByID", ctx, userID).Return(nil, assert.AnError)

	req := &pb.GetUserRequest{Id: int32(userID)}

	// Act
	resp, err := service.GetUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Пользователь не найден")

	mockRepo.AssertExpectations(t)
}
