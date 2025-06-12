# Изображения для документации

Эта папка содержит изображения для README.md файла.

## Как добавить реальные изображения:

### 1. graylog-dashboard.png
- Откройте Graylog: http://localhost:9000
- Войдите (admin/admin)
- Сделайте скриншот дашборда с логами приложения
- Сохраните как `graylog-dashboard.png`

### 2. grafana.png
- Откройте Grafana: http://localhost:3001
- Войдите (admin/admin)
- Настройте Prometheus data source
- Сделайте скриншот дашборда с метриками
- Сохраните как `grafana.png`

### 3. architecture-diagram.png
Создайте диаграмму архитектуры, показывающую:
```
[React App] -> [HTTP Gateway] -> [gRPC Server] -> [PostgreSQL]
                     |
                [Monitoring]
                     |
        [Prometheus] [Grafana] [Graylog]
```

### 4. web-app-screenshot.png
- Откройте веб-приложение: http://localhost:3000
- Сделайте скриншот интерфейса с пользователями
- Сохраните как `web-app-screenshot.png`

## Рекомендуемые размеры:
- Ширина: 800-1200px
- Формат: PNG или JPG
- Качество: высокое, но оптимизированное для веба

## Инструменты для создания диаграмм:
- [Draw.io](https://draw.io)
- [Lucidchart](https://lucidchart.com)
- [Mermaid](https://mermaid.js.org)
- [PlantUML](https://plantuml.com) 