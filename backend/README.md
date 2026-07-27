# Бэкенд «Мы Друзья» (API-заглушка на Litestar)

Небольшой ASGI-сервер на [Litestar](https://litestar.dev) (Python). Пока это
**заглушка**: данные лежат в `app/data/*.py` (те же объекты, что раньше были
захардкожены во фронте) и отдаются как JSON. Контракты ручек —
в [`../docs/api-contracts.md`](../docs/api-contracts.md).

## Установка и запуск

```bash
cd backend
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
litestar --app app.main:app run --port 3001 --reload
```

Сервер поднимется на `http://localhost:3001`. Проверка: `http://localhost:3001/api`.

## Ручки

| Метод | Путь                 | Отдаёт            |
|-------|----------------------|-------------------|
| GET   | `/api/news`          | Список новостей   |
| GET   | `/api/news/{id}`     | Новость целиком   |
| GET   | `/api/classes`       | Список занятий    |
| GET   | `/api/classes/{id}`  | Занятие целиком   |
| GET   | `/api/services`      | Список услуг      |
| GET   | `/api/services/{id}` | Услуга целиком    |
| GET   | `/api/staff`         | Педагоги          |
| GET   | `/api/documents`     | Документы         |
| GET   | `/api/gallery`       | Галерея           |

## Как это связано с фронтом

Фронт (Vite) в дев-режиме проксирует `/api/*` на этот сервер (см.
`frontend/vite.config.ts`), поэтому CORS в разработке не мешает. Данные можно
наполнять здесь — во фронте менять ничего не нужно.

Когда появится настоящая база данных, замените источники данных в `app/main.py`
(импорты из `app/data/…`) на запросы к БД, сохранив формат ответов. Ключи в JSON
намеренно в camelCase, чтобы совпадать с типами фронтенда
(`frontend/src/data/types.ts`).
