# Go backend «Мы Друзья»

Основной API на Go, PostgreSQL и MinIO. Litestar-заглушка для изолированной
разработки находится в [`../backend-stub`](../backend-stub).

JSON-контракты не меняются: они зафиксированы в
[`../docs/api-contracts.md`](../docs/api-contracts.md).

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
├── docker-compose.yml       # PostgreSQL и MinIO для разработки
├── Makefile
└── .env.example
```

Реализованы публичные и административные ручки, PostgreSQL-хранилище и MinIO.

## Локальный запуск

1. Создайте `backend-go/.env` на основе `.env.example` и задайте только секреты.
2. Экспортируйте переменные из файла: `set -a; source .env; set +a`.
3. Запустите инфраструктуру и схему: `make up && make migrate`.
4. Запустите API: `make run`.

Профиль `development` включён по умолчанию. Для production используется
`APP_ENV=production`. Несекретные значения заданы в `internal/config` и могут
быть переопределены через env.
