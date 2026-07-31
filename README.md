# we-are-friends

Сайт детского центра «Мы Друзья».

- [`frontend/`](frontend) — Vue-клиент.
- [`backend-stub/`](backend-stub) — временный Litestar API с контентом в Python.
- [`backend-go/`](backend-go) — новый Go API на PostgreSQL и MinIO; его план и
  этапы переноса описаны в [`docs/backend-go-plan.md`](docs/backend-go-plan.md).

## Запуск в Docker

1. Скопируйте `.env.example` в `.env` и задайте сильные `ADMIN_PASSWORD`,
   `POSTGRES_PASSWORD` и `MINIO_ROOT_PASSWORD`.
2. Запустите инфраструктуру и приложение: `docker-compose up --build -d`.
3. Один раз импортируйте имеющийся контент: `docker-compose --profile seed run --rm seed`.
   После импорта временный сервис-заглушку можно остановить:
   `docker-compose --profile seed stop stub`.

Сайт, включая `/admin`, доступен на `http://localhost:8088`. Nginx раздаёт
собранный Vue-клиент, проксирует `/api` в Go и `/media` в MinIO, а также
возвращает `index.html` для всех клиентских маршрутов. Vite не запускается как
сервер и не участвует в маршрутизации: он остаётся только одноразовым
компилятором `.vue`-файлов в build-слое Docker.
