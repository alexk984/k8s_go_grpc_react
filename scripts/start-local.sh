#!/bin/bash

echo "🚀 Запуск K8s Go gRPC React приложения локально"

# Проверяем зависимости
command -v docker >/dev/null 2>&1 || { echo "❌ Docker не установлен"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose не установлен"; exit 1; }

# Генерируем protobuf файлы если их нет
if [ ! -f "proto/user.pb.go" ]; then
    echo "📝 Генерируем protobuf файлы..."
    export PATH="$PATH:$(go env GOPATH)/bin"
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto
fi

# Запускаем все сервисы
echo "🐳 Запускаем все сервисы с Docker Compose..."
docker-compose up --build -d

echo "⏳ Ждем инициализации сервисов..."
sleep 10

echo "✅ Сервисы запущены!"
echo ""
echo "🌐 Доступные сервисы:"
echo "   React App:    http://localhost:3000"
echo "   API Gateway:  http://localhost:8081"
echo "   Grafana:      http://localhost:3001 (admin/admin)"
echo "   Prometheus:   http://localhost:9091"
echo ""
echo "📊 Проверка состояния:"
curl -s http://localhost:8081/health && echo " ✅ HTTP Gateway работает"
curl -s http://localhost:9091/-/healthy && echo " ✅ Prometheus работает"
echo ""
echo "🔍 Просмотр логов: docker-compose logs -f"
echo "🛑 Остановка: docker-compose down" 