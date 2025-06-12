# Примеры использования JWT аутентификации

## Запуск сервера

```bash
# Установка переменных окружения
export DATABASE_URL="postgres://user:password@localhost:5432/k8s_grpc_db?sslmode=disable"
export JWT_SECRET="your-secret-key"
export GRPC_PORT="8080"
export HTTP_PORT="8081"

# Запуск сервера
go run cmd/server/main.go
```

## HTTP API (через gRPC-Gateway)

### 1. Регистрация нового пользователя

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

Ответ:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "is_active": true,
    "created_at": 1640995200
  },
  "message": "Пользователь успешно зарегистрирован"
}
```

### 2. Вход в систему

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

Ответ:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "is_active": true,
    "created_at": 1640995200
  },
  "message": "Успешный вход в систему"
}
```

### 3. Получение пользователя (требует аутентификации)

```bash
# Сохраните токен из предыдущего запроса
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8081/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Получение списка пользователей (требует аутентификации)

```bash
curl -X GET http://localhost:8081/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Создание пользователя (требует аутентификации)

```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane@example.com",
    "password": "password456",
    "role": "admin"
  }'
```

## gRPC API (прямое подключение)

### Использование grpcurl

```bash
# Установка grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Регистрация
grpcurl -plaintext -d '{
  "name": "John Doe",
  "email": "john@example.com", 
  "password": "password123"
}' localhost:8080 user.UserService/Register

# Вход
grpcurl -plaintext -d '{
  "email": "john@example.com",
  "password": "password123"
}' localhost:8080 user.UserService/Login

# Получение пользователя (с токеном)
grpcurl -plaintext \
  -H "authorization: Bearer $TOKEN" \
  -d '{"id": 1}' \
  localhost:8080 user.UserService/GetUser
```

## Структура JWT токена

Токен содержит следующие claims:
- `user_id`: ID пользователя
- `email`: Email пользователя  
- `role`: Роль пользователя (user, admin)
- `exp`: Время истечения токена (по умолчанию 24 часа)
- `iat`: Время создания токена

## Middleware аутентификации

### Публичные методы (не требуют токена):
- `Register` - регистрация
- `Login` - вход в систему

### Защищенные методы (требуют токен):
- `GetUser` - получение пользователя
- `CreateUser` - создание пользователя
- `ListUsers` - список пользователей

## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DATABASE_URL` | URL подключения к PostgreSQL | `postgres://user:password@localhost:5432/k8s_grpc_db?sslmode=disable` |
| `JWT_SECRET` | Секретный ключ для JWT | `your-secret-key` |
| `GRPC_PORT` | Порт gRPC сервера | `8080` |
| `HTTP_PORT` | Порт HTTP сервера (gRPC-Gateway) | `8081` |

## Мониторинг

### Метрики Prometheus

```bash
# Просмотр метрик
curl http://localhost:8081/metrics
```

### Health Check

```bash
# Проверка состояния сервера
curl http://localhost:8081/health
``` 