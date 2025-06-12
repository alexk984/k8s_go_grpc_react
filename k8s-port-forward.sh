#!/bin/bash

# Настройка port-forward для всех сервисов K8s приложения

echo "🚀 Настройка port-forward для K8s сервисов..."

# Веб-приложение
echo "📱 Запуск port-forward для веб-приложения на порту 3000..."
kubectl port-forward service/k8s-grpc-app-web-app 3000:80 &

# HTTP Gateway
echo "🔗 Запуск port-forward для HTTP Gateway на порту 8081..."
kubectl port-forward service/k8s-grpc-app-http-gateway 8081:8081 &

# Grafana
echo "📊 Запуск port-forward для Grafana на порту 3001..."
kubectl port-forward service/k8s-grpc-app-grafana 3001:3000 &

# Prometheus
echo "🔍 Запуск port-forward для Prometheus на порту 9091..."
kubectl port-forward service/k8s-grpc-app-prometheus 9091:9090 &

# Graylog
echo "📝 Запуск port-forward для Graylog на порту 9000..."
kubectl port-forward service/k8s-grpc-app-graylog 9000:9000 &

# PostgreSQL
echo "🗄️  Запуск port-forward для PostgreSQL на порту 5432..."
kubectl port-forward service/k8s-grpc-app-postgres 5432:5432 &

echo ""
echo "✅ Все port-forward настроены!"
echo ""
echo "🌐 Доступные сервисы:"
echo "  • Веб-приложение: http://localhost:3000"
echo "  • HTTP API:       http://localhost:8081"
echo "  • Grafana:        http://localhost:3001 (admin/admin)"
echo "  • Prometheus:     http://localhost:9091"
echo "  • Graylog:        http://localhost:9000 (admin/admin)"
echo "  • PostgreSQL:     localhost:5432 (postgres/postgres)"
echo ""
echo "⏹️  Для остановки: Ctrl+C"
echo "🔄 Ожидание завершения..."

# Ждем завершения всех фоновых процессов
wait 