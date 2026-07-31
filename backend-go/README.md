# Go backend «Мы Друзья»

Это каркас будущего API на Go, PostgreSQL и MinIO. Существующая рабочая
Litestar-заглушка находится в [`../backend-stub`](../backend-stub) и остаётся
источником данных до завершения переноса.

JSON-контракты не меняются: они зафиксированы в
[`../docs/api-contracts.md`](../docs/api-contracts.md). Подробный порядок работ,
схема PostgreSQL и принятые архитектурные решения — в
[`../docs/backend-go-plan.md`](../docs/backend-go-plan.md).

## Структура

```
backend-go/
├── cmd/api/                 # исполняемый HTTP-сервис
├── internal/
│   ├── config/              # переменные окружения и их проверка
│   ├── httpapi/             # роутер, DTO, хендлеры и HTTP-ошибки
│   ├── media/               # MinIO и URL для объектов
│   └── store/               # SQL-запросы и генерируемый sqlc код
├── migrations/              # миграции goose
├── seed/                    # импорт контента и фотографий из заглушки
├── docker-compose.yml       # PostgreSQL и MinIO для разработки
├── Makefile
└── .env.example
```

Реализованы публичные ручки из контракта, PostgreSQL-хранилище, MinIO-URL и
импорт текущего контента. Админские CRUD-ручки остаются отдельным следующим
этапом плана.

## Локальный запуск

1. Создайте `backend-go/.env` на основе `.env.example`. Для compose-разработки
   используйте `DATABASE_URL=postgres://we_are_friends:we_are_friends@localhost:15432/we_are_friends?sslmode=disable`,
   `MINIO_ENDPOINT=localhost:9000`, `MINIO_ACCESS_KEY=minioadmin`,
   `MINIO_SECRET_KEY=minioadmin`, `MINIO_BUCKET=media`,
   `MEDIA_PUBLIC_BASE_URL=http://localhost:9000/media`, `PORT=8080`.
2. Экспортируйте переменные из файла: `set -a; source .env; set +a`.
3. Запустите инфраструктуру и схему: `make up && make migrate`.
4. В отдельном терминале запустите заглушку на порту 3001, затем импортируйте
   её контент: `make seed`.
5. Запустите API: `make run`.

Сервис слушает `PORT` (по умолчанию `8080`). Заглушка нужна только во время
первоначального импорта; после него Vite ходит к Go API.
