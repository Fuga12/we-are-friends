# Мы Друзья

Сайт детского центра: Vue, Go, PostgreSQL и MinIO.

## Запуск

Нужны Docker и Docker Compose.

```bash
cp .env.example .env
```

Заполните в `.env` три пароля:

```env
POSTGRES_PASSWORD=сложный-пароль
MINIO_ROOT_PASSWORD=сложный-пароль
ADMIN_PASSWORD=пароль-админки
```

Запустите сайт:

```bash
docker-compose up -d --build
```

- Сайт: `http://localhost:8088`
- Админка: `http://localhost:8088/admin`
- Логин админки: `admin`

## Выкладка на сервер

```bash
git clone <адрес-репозитория> we-are-friends
cd we-are-friends
cp .env.example .env
nano .env
docker-compose up -d --build
```

Направьте домен через Nginx или другой reverse proxy на `127.0.0.1:8088`.

## Обновление

```bash
git pull
docker-compose up -d --build
```

## Полезные команды

```bash
docker-compose ps
docker-compose logs -f
docker-compose down
```

Данные хранятся в папке `data/`. Не удаляйте её и добавьте в резервное копирование.
