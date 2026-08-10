# Backend

Go API сайта «Мы Друзья». Использует PostgreSQL и MinIO.

Обычный запуск выполняется из корня проекта:

```bash
docker-compose up -d --build
```

Проверка без Docker:

```bash
GOCACHE=/tmp/we-are-friends-go-cache go test ./...
```

Миграции применяются автоматически при запуске API.
