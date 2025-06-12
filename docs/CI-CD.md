# CI/CD Pipeline

Этот проект использует GitHub Actions для автоматизации процессов сборки, тестирования и развертывания.

## Структура Pipeline

### 1. Основной CI/CD Pipeline (`.github/workflows/ci-cd.yml`)

Запускается при:
- Push в ветки `main` и `develop`
- Pull Request в ветку `main`

**Этапы:**
1. **Test** - Запуск тестов Go и React
2. **Build** - Сборка и публикация Docker образов
3. **Security Scan** - Сканирование безопасности с Trivy
4. **Deploy Staging** - Развертывание в staging (ветка `develop`)
5. **Deploy Production** - Развертывание в production (ветка `main`)
6. **Notify** - Уведомления в Slack

### 2. PR Checks (`.github/workflows/pr-checks.yml`)

Быстрые проверки для Pull Request:
- Линтинг и форматирование кода
- Unit тесты
- Проверка сборки Docker образов

### 3. Release Pipeline (`.github/workflows/release.yml`)

Запускается при создании тегов `v*`:
- Создание GitHub Release
- Сборка и публикация релизных образов
- Обновление Helm чартов
- Упаковка и публикация Helm чартов

## Настройка

### 1. GitHub Secrets

Необходимо настроить следующие секреты в GitHub:

```bash
# Kubernetes конфигурации
KUBE_CONFIG_STAGING      # Base64 encoded kubeconfig для staging
KUBE_CONFIG_PRODUCTION   # Base64 encoded kubeconfig для production

# Slack уведомления (опционально)
SLACK_WEBHOOK           # Webhook URL для Slack уведомлений
```

### 2. GitHub Environments

Создайте environments в GitHub:
- `staging` - для staging развертывания
- `production` - для production развертывания (с защитой)

### 3. Container Registry

Pipeline использует GitHub Container Registry (ghcr.io) для хранения Docker образов.
Убедитесь, что у репозитория есть права на запись в packages.

## Docker Images

Pipeline создает следующие образы:

```
ghcr.io/[username]/[repo]/grpc-server:latest
ghcr.io/[username]/[repo]/http-gateway:latest
ghcr.io/[username]/[repo]/web-app:latest
```

Теги:
- `latest` - для main ветки
- `[branch-name]-[sha]` - для других веток
- `v[version]` - для релизов

## Тестирование

### Go тесты
```bash
go test -v ./...
```

### React тесты
```bash
cd web
npm test -- --coverage --watchAll=false
```

### Линтинг
```bash
# Go
golangci-lint run

# React
cd web
npm run lint
npm run type-check
npm run format:check
```

## Развертывание

### Staging
Автоматически развертывается при push в ветку `develop`.

### Production
Автоматически развертывается при push в ветку `main`.

### Manual Deploy
```bash
# Staging
helm upgrade --install k8s-grpc-app-staging ./helm/k8s-grpc-app \
  --namespace staging \
  --create-namespace \
  --set environment=staging

# Production
helm upgrade --install k8s-grpc-app ./helm/k8s-grpc-app \
  --namespace production \
  --create-namespace \
  --set environment=production
```

## Мониторинг

Pipeline включает:
- Метрики времени выполнения
- Покрытие кода (Codecov)
- Сканирование безопасности (Trivy)
- Уведомления в Slack

## Создание Release

1. Создайте и push тег:
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. GitHub Actions автоматически:
   - Создаст GitHub Release
   - Соберет и опубликует образы с версией
   - Обновит Helm чарты
   - Упакует Helm чарт для скачивания

## Troubleshooting

### Ошибки сборки
- Проверьте логи в GitHub Actions
- Убедитесь, что все тесты проходят локально
- Проверьте форматирование кода

### Ошибки развертывания
- Проверьте корректность kubeconfig
- Убедитесь, что namespace существует
- Проверьте права доступа к кластеру

### Ошибки Docker Registry
- Проверьте права доступа к GitHub Packages
- Убедитесь, что токен имеет scope `write:packages`

## Локальная разработка

Для локального тестирования pipeline:

```bash
# Установка act (GitHub Actions локально)
brew install act

# Запуск тестов локально
act -j test

# Запуск с секретами
act -j test --secret-file .secrets
```

## Best Practices

1. **Всегда создавайте PR** для изменений
2. **Пишите тесты** для нового кода
3. **Следуйте соглашениям** по именованию коммитов
4. **Проверяйте линтеры** перед push
5. **Используйте семантическое версионирование** для релизов 