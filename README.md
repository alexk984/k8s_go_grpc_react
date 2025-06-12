# 🚀 K8s Go gRPC React приложение

[![Go](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18-blue.svg)](https://reactjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue.svg)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://www.postgresql.org/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-green.svg)](https://kubernetes.io/)
[![Docker](https://img.shields.io/badge/Docker-Compose-blue.svg)](https://www.docker.com/)
[![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-green.svg)](https://github.com/features/actions)
[![Tests](https://img.shields.io/badge/Tests-Passing-green.svg)](#)

Микросервисная архитектура с полным стеком мониторинга для Kubernetes развертывания.

![Graylog Dashboard](docs/images/graylog-dashboard.png)
*Централизованное логирование в Graylog с структурированными логами от всех микросервисов*

## ✨ Особенности

- 🔄 **Микросервисная архитектура** с gRPC и HTTP Gateway
- 🗄️ **PostgreSQL** с GORM ORM для надежного хранения данных
- 🎯 **TypeScript** фронтенд с полной типизацией
- 📊 **Полный стек мониторинга**: Prometheus + Grafana + Graylog
- 🚀 **Kubernetes-ready** с Helm charts
- 🐳 **Docker Compose** для локальной разработки
- 📝 **Структурированное логирование** в формате GELF
- 🔍 **Автоматические метрики** для всех API endpoints
- 🛡️ **Production-ready** конфигурация
- 🔄 **CI/CD Pipeline** с GitHub Actions
- 🧪 **Автоматическое тестирование** Go и React кода

## 🏗️ Архитектура

![Architecture Diagram](docs/images/architecture-diagram.png)
*Диаграмма архитектуры микросервисного приложения*

- **Go gRPC Server** - основной сервис для управления пользователями с PostgreSQL
- **HTTP Gateway** - мост между веб-фронтендом и gRPC backend
- **Веб-приложение** - современный React TypeScript интерфейс для взаимодействия с API
- **PostgreSQL** - надежная реляционная база данных для хранения пользователей
- **Prometheus** - сбор метрик 
- **Grafana** - визуализация метрик и дашборды
- **Graylog** - централизованное логирование

## 🔧 Технологии

- **Backend**: Go 1.23, gRPC, Protocol Buffers, GORM
- **База данных**: PostgreSQL 15
- **Frontend**: React 18, TypeScript, CSS3
- **Мониторинг**: Prometheus, Grafana, Graylog
- **Контейнеризация**: Docker, Docker Compose
- **Оркестрация**: Kubernetes, Helm
- **CI/CD**: GitHub Actions, Automated Testing, Security Scanning

## ⚡ Быстрый запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/your-username/k8s_go_grpc_react.git
cd k8s_go_grpc_react
```

### 2. Запуск с Docker Compose (рекомендуется)
```bash
# Запуск всех сервисов
docker-compose up -d

# Проверка статуса
docker-compose ps
```

### 3. Доступ к сервисам

#### Основные сервисы:
- **Веб-приложение**: http://localhost:3000 
- **HTTP API Gateway**: http://localhost:8081
- **gRPC Server**: localhost:8080

#### Мониторинг:
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9091
- **Graylog**: http://localhost:9000 (admin/admin) - только в Kubernetes

#### База данных:
- **PostgreSQL**: localhost:5432 (postgres/postgres) - только в Kubernetes

### 4. Тестирование API

```bash
# Проверка здоровья
curl http://localhost:8081/health

# Регистрация нового пользователя
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Иван Иванов","email":"ivan@example.com","password":"password123"}'

# Вход в систему
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ivan@example.com","password":"password123"}'

# Получение списка пользователей (с JWT токеном)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8081/api/v1/users

# Получение пользователя по ID
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8081/api/v1/users/1
```

## 🌐 Веб-интерфейс

![Web Application](docs/images/web-app-screenshot.png)
*React TypeScript приложение для управления пользователями с JWT аутентификацией*

Откройте http://localhost:3000 в браузере для доступа к современному React TypeScript интерфейсу с возможностями:

- ✅ **JWT аутентификация** - регистрация и вход в систему
- ✅ **Управление пользователями** с защитой ролями
- ✅ **Автоматический logout** при истечении токена
- ✅ **Защищенные маршруты** для авторизованных пользователей
- ✅ **Профиль пользователя** с информацией о сессии
- ✅ **Просмотр списка пользователей** в реальном времени
- ✅ **Ссылки на мониторинг** (Grafana, Prometheus, Graylog)
- ✅ **Типизированная обработка ошибок** с TypeScript
- ✅ **Адаптивный дизайн** для всех устройств
- ✅ **Автоматическое обновление токенов** и обработка 401 ошибок

## 📊 Мониторинг и метрики

### Grafana Dashboards
![Grafana Dashboard](docs/images/grafana.png)
*Дашборд Grafana с метриками производительности и мониторингом API*

1. Откройте http://localhost:3001
2. Войдите (admin/admin)
3. Добавьте Prometheus data source: http://prometheus:9090
4. Импортируйте дашборды из `monitoring/grafana/dashboards/`

### Prometheus Metrics
- API запросы и ответы
- Время выполнения
- Количество ошибок
- Системные метрики

### Graylog Logs
![Graylog Logs](docs/images/graylog-dashboard.png)
*Структурированные логи в Graylog с фильтрацией по компонентам*

- Централизованные логи всех сервисов
- Структурированное логирование в формате GELF
- Мощный поиск и фильтрация
- Веб-интерфейс для анализа логов

## 🐳 Развертывание в Kubernetes

### Предварительные требования
- Kubernetes кластер (minikube, Docker Desktop или облачный)
- Helm 3+
- kubectl

### Установка через Helm
```bash
# 1. Запуск minikube (если используется)
minikube start

# 2. Загрузка Docker образов в minikube (если используется)
eval $(minikube docker-env)
docker build -t k8s-grpc-app-grpc-server:latest -f Dockerfile.server .
docker build -t k8s-grpc-app-http-gateway:latest -f Dockerfile.gateway .
docker build -t k8s-grpc-app-web-app:latest -f web/Dockerfile web/

# 3. Установка приложения с мониторингом в namespace default
helm install k8s-grpc-app ./helm/k8s-grpc-app

# 4. Проверка статуса подов
kubectl get pods | grep k8s-grpc-app
kubectl get services | grep k8s-grpc-app
```

### Доступ к сервисам
```bash
# Запуск port-forward для всех сервисов (рекомендуется)
chmod +x k8s-port-forward.sh
./k8s-port-forward.sh

# Доступные сервисы:
# • Веб-приложение: http://localhost:3000
# • HTTP API:       http://localhost:8081
# • Grafana:        http://localhost:3001 (admin/admin)
# • Prometheus:     http://localhost:9091
# • Graylog:        http://localhost:9000 (admin/admin)

# Или по отдельности:
kubectl port-forward service/k8s-grpc-app-web-app 3000:80 &
kubectl port-forward service/k8s-grpc-app-http-gateway 8081:8081 &
kubectl port-forward service/k8s-grpc-app-grafana 3001:3000 &
kubectl port-forward service/k8s-grpc-app-prometheus 9091:9090 &
kubectl port-forward service/k8s-grpc-app-graylog 9000:9000 &
```

### Проверка работы
```bash
# Тест API (сначала войти и получить токен)
curl http://localhost:8081/health

# Регистрация пользователя
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"K8s User","email":"k8s@example.com","password":"password123"}'

# Вход в систему
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"k8s@example.com","password":"password123"}'

# Получение списка пользователей (с JWT токеном)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8081/api/v1/users
```

### Удаление
```bash
# Удаление приложения
helm uninstall k8s-grpc-app

# Удаление persistent volumes (опционально)
kubectl delete pvc --all
```

## 🛠️ Разработка

### Структура проекта
```
.
├── cmd/                    # Точки входа приложения
│   ├── server/            # gRPC сервер
│   └── gateway/           # HTTP шлюз
├── internal/              # Внутренняя логика
│   ├── service/          # Бизнес логика
│   └── logger/           # Graylog интеграция
├── proto/                 # Protocol Buffers схемы
├── web/                   # React TypeScript приложение
│   ├── src/
│   │   ├── types/        # TypeScript типы
│   │   ├── services/     # API сервисы
│   │   ├── App.tsx       # Главный компонент
│   │   └── index.tsx     # Точка входа
│   ├── public/           # Статические файлы
│   └── tsconfig.json     # TypeScript конфигурация
├── monitoring/            # Конфигурации мониторинга
├── helm/                  # Helm чарты
└── docker-compose.yml     # Локальная разработка
```

### Генерация protobuf
```bash
# Генерация Go кода из proto файлов
./scripts/generate-proto.sh
```

### Локальная разработка

#### Backend (Go)
```bash
# Запуск только gRPC сервера
go run cmd/server/main.go

# Запуск только HTTP gateway
go run cmd/gateway/main.go
```

#### Frontend (React TypeScript)
```bash
# Переход в папку веб-приложения
cd web

# Установка зависимостей
npm install

# Запуск в режиме разработки
npm start

# Сборка для продакшена
npm run build

# Проверка типов TypeScript
npx tsc --noEmit
```

## 🔍 Отладка

### Просмотр логов
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f grpc-server
docker-compose logs -f http-gateway
```

### Остановка сервисов
```bash
# Остановка всех контейнеров
docker-compose down

# Остановка с удалением volumes
docker-compose down -v
```

## 📋 API Endpoints

### Публичные endpoints
| Method | Endpoint | Описание |
|--------|----------|----------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Регистрация пользователя |
| POST | `/api/v1/auth/login` | Вход в систему |

### Защищенные endpoints (требуют JWT токен)
| Method | Endpoint | Описание |
|--------|----------|----------|
| GET | `/api/v1/users` | Получить всех пользователей |
| GET | `/api/v1/users/{id}` | Получить пользователя по ID |
| POST | `/api/v1/users` | Создать пользователя |

### Примеры запросов
```json
// POST /api/v1/auth/register
{
  "name": "Иван Иванов",
  "email": "ivan@example.com",
  "password": "password123"
}

// Response
{
  "user": {
    "id": 1,
    "name": "Иван Иванов", 
    "email": "ivan@example.com",
    "role": "user",
    "created_at": 1749598180
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Пользователь успешно зарегистрирован"
}

// POST /api/v1/auth/login
{
  "email": "ivan@example.com",
  "password": "password123"
}

// Response
{
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@example.com",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Авторизация успешна"
}

// GET /api/v1/users (с JWT токеном)
// Headers: Authorization: Bearer YOUR_JWT_TOKEN
// Response
{
  "users": [
    {
      "id": 1,
      "name": "Иван Иванов",
      "email": "ivan@example.com",
      "role": "user",
      "created_at": 1749598180
    }
  ]
}
```

## 🔄 CI/CD Pipeline

Проект использует GitHub Actions для автоматизации процессов разработки:

### Автоматические проверки
- ✅ **Линтинг кода** (Go и TypeScript)
- ✅ **Unit тесты** для всех компонентов
- ✅ **Проверка типов** TypeScript
- ✅ **Сканирование безопасности** с Trivy
- ✅ **Проверка сборки** Docker образов

### Автоматическое развертывание
- 🚀 **Staging** - автоматически при push в `develop`
- 🚀 **Production** - автоматически при push в `main`
- 📦 **Docker образы** публикуются в GitHub Container Registry
- 🏷️ **Релизы** создаются автоматически при создании тегов

### Workflow файлы
- `.github/workflows/ci-cd.yml` - основной pipeline
- `.github/workflows/pr-checks.yml` - проверки для PR
- `.github/workflows/release.yml` - создание релизов

### Настройка CI/CD
Подробная документация по настройке и использованию CI/CD pipeline доступна в [docs/CI-CD.md](docs/CI-CD.md).

### Создание релиза
```bash
# Создание и push тега для релиза
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions автоматически:
# - Создаст GitHub Release
# - Соберет и опубликует Docker образы
# - Обновит Helm чарты
```

## 🧪 Тестирование

### Запуск тестов локально
```bash
# Go тесты
go test -v ./...

# React тесты
cd web
npm test -- --coverage --watchAll=false

# Линтинг
golangci-lint run
cd web && npm run lint
```

### Покрытие кода
Отчеты о покрытии кода автоматически генерируются в CI/CD pipeline и отправляются в Codecov.

# Kubernetes Go gRPC React Application

[![CI/CD Pipeline](https://github.com/username/k8s_go_grpc_react/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/username/k8s_go_grpc_react/actions/workflows/ci-cd.yml)
[![PR Checks](https://github.com/username/k8s_go_grpc_react/actions/workflows/pr-checks.yml/badge.svg)](https://github.com/username/k8s_go_grpc_react/actions/workflows/pr-checks.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/username/k8s_go_grpc_react)](https://goreportcard.com/report/github.com/username/k8s_go_grpc_react)

Современное веб-приложение с микросервисной архитектурой, использующее Go, gRPC, React и Kubernetes.

## 🚀 Возможности

- **gRPC API** с автоматической генерацией HTTP REST API через gRPC-Gateway
- **JWT аутентификация** с middleware для защиты эндпоинтов
- **PostgreSQL** интеграция с GORM ORM
- **React фронтенд** с TypeScript и современным UI
- **Kubernetes** деплой с Helm чартами
- **CI/CD пайплайн** с GitHub Actions
- **Мониторинг** с Prometheus метриками
- **Автоматическое тестирование** Go и React компонентов
- **Безопасность** сканирование с Trivy
- **Линтинг и форматирование** кода

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │   Go gRPC       │    │   PostgreSQL    │
│   (Frontend)    │◄──►│   Server        │◄──►│   Database      │
│   Port: 3000    │    │   Port: 8080    │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   HTTP Gateway  │
                       │   Port: 8081    │
                       └─────────────────┘
```

## 🔐 Аутентификация

Приложение использует JWT (JSON Web Tokens) для аутентификации:

### Публичные эндпоинты:
- `POST /api/v1/auth/register` - регистрация пользователя
- `POST /api/v1/auth/login` - вход в систему

### Защищенные эндпоинты (требуют JWT токен):
- `GET /api/v1/users/{id}` - получение пользователя
- `GET /api/v1/users` - список пользователей
- `POST /api/v1/users` - создание пользователя

Подробные примеры использования см. в [examples/auth_example.md](examples/auth_example.md)

## 🛠️ Технологический стек

### Backend
- **Go 1.21+** - основной язык программирования
- **gRPC** - межсервисное взаимодействие
- **gRPC-Gateway** - автоматическая генерация REST API
- **PostgreSQL** - основная база данных
- **GORM** - ORM для работы с базой данных
- **JWT** - аутентификация и авторизация
- **Prometheus** - метрики и мониторинг

### Frontend
- **React 18** - пользовательский интерфейс
- **TypeScript** - типизированный JavaScript
- **Material-UI** - компоненты интерфейса
- **Axios** - HTTP клиент

### DevOps
- **Kubernetes** - оркестрация контейнеров
- **Helm** - управление Kubernetes манифестами
- **Docker** - контейнеризация
- **GitHub Actions** - CI/CD пайплайн
- **Trivy** - сканирование безопасности

## 🚀 Быстрый старт

### Предварительные требования

- Go 1.21+
- Node.js 18+
- PostgreSQL 13+
- Docker (опционально)
- Kubernetes (для деплоя)

### Локальная разработка

1. **Клонирование репозитория**
```bash
git clone https://github.com/username/k8s_go_grpc_react.git
cd k8s_go_grpc_react
```

2. **Настройка базы данных**
```bash
# Создание базы данных PostgreSQL
createdb k8s_grpc_db

# Установка переменных окружения
export DATABASE_URL="postgres://user:password@localhost:5432/k8s_grpc_db?sslmode=disable"
export JWT_SECRET="your-secret-key"
```

3. **Запуск backend сервера**
```bash
# Установка зависимостей
go mod download

# Генерация protobuf файлов
./scripts/generate_proto.sh

# Запуск сервера
go run cmd/server/main.go
```

4. **Запуск frontend приложения**
```bash
cd web
npm install
npm start
```

5. **Тестирование API**
```bash
# Регистрация пользователя
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'

# Вход в систему
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password123"}'
```

## 🧪 Тестирование

### Backend тесты
```bash
# Запуск всех тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Тесты конкретного пакета
go test ./internal/service/...
```

### Frontend тесты
```bash
cd web
npm test
```

### Линтинг
```bash
# Go линтинг
golangci-lint run

# TypeScript линтинг
cd web && npm run lint
```

## 📦 Деплой

### Docker
```bash
# Сборка образов
docker build -t k8s-grpc-app:latest .
docker build -t k8s-grpc-web:latest ./web

# Запуск с docker-compose
docker-compose up -d
```

### Kubernetes
```bash
# Деплой с Helm
helm install k8s-grpc-app ./helm/k8s-grpc-app \
  --set database.url="postgres://user:password@postgres:5432/k8s_grpc_db"
```

## 📊 Мониторинг

### Метрики Prometheus
- `grpc_requests_total` - общее количество gRPC запросов
- `grpc_request_duration_seconds` - время выполнения запросов
- `users_total` - общее количество пользователей

### Эндпоинты мониторинга
- `GET /metrics` - метрики Prometheus
- `GET /health` - проверка состояния сервера

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DATABASE_URL` | URL подключения к PostgreSQL | `postgres://user:password@localhost:5432/k8s_grpc_db?sslmode=disable` |
| `JWT_SECRET` | Секретный ключ для JWT | `your-secret-key` |
| `GRPC_PORT` | Порт gRPC сервера | `8080` |
| `HTTP_PORT` | Порт HTTP сервера | `8081` |

## 📚 Документация

- [Примеры использования API](examples/auth_example.md)
- [CI/CD Pipeline](docs/CI-CD.md)
- [Архитектура приложения](docs/ARCHITECTURE.md)

## 🤝 Участие в разработке

1. Форкните репозиторий
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 🙏 Благодарности

- [gRPC](https://grpc.io/) - высокопроизводительный RPC фреймворк
- [React](https://reactjs.org/) - библиотека для создания пользовательских интерфейсов
- [Kubernetes](https://kubernetes.io/) - система оркестрации контейнеров
- [PostgreSQL](https://www.postgresql.org/) - мощная реляционная база данных
