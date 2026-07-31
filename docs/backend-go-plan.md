# План настоящего бэкенда: Go + PostgreSQL + MinIO

Документ описывает, чем и как заменить текущую заглушку на Litestar
(`backend-stub/`, данные в `app/data/*.py`). Контракты ручек уже зафиксированы в
[api-contracts.md](api-contracts.md) — **новый бэкенд обязан отдавать тот же
JSON**, тогда фронтенд не придётся переделывать вообще.

---

## 1. Стек и почему именно он

Правило выбора: берём то, что чаще всего встречается в чужих проектах и
туториалах — так проще гуглить ошибки и читать примеры.

| Задача | Выбор | Почему |
|--------|-------|--------|
| HTTP-сервер | стандартный `net/http` | Полноценный сервер в стандартной библиотеке. С Go 1.22 `http.ServeMux` умеет методы и параметры пути: `GET /api/classes/{id}` |
| Роутер / middleware | [`go-chi/chi/v5`](https://github.com/go-chi/chi) | Тонкая надстройка над `net/http`: группы маршрутов, готовые middleware (логи, CORS, recover, таймауты). Хендлеры — обычные `http.HandlerFunc`, то есть навык переносится куда угодно |
| Драйвер Postgres | [`jackc/pgx/v5`](https://github.com/jackc/pgx) | Де-факто стандарт для Postgres в Go. Быстрее `database/sql`, поддерживает массивы, `jsonb`, пул соединений |
| Запросы к БД | [`sqlc`](https://sqlc.dev) | Пишешь обычный SQL — генератор делает из него типизированные Go-функции. Ошибки в SQL видны при генерации, а не в рантайме |
| Миграции | [`pressly/goose`](https://github.com/pressly/goose) | Миграции — обычные `.sql` файлы с `-- +goose Up/Down`. Умеет вшиваться в бинарник и накатываться при старте |
| Файлы (S3) | [`minio-go/v7`](https://github.com/minio/minio-go) | Официальный SDK MinIO, он же S3-совместимый: при переезде на Yandex Object Storage / AWS меняется только endpoint |
| Логи | `log/slog` (стандартная библиотека) | С Go 1.21 структурированные логи есть из коробки, внешний логгер не нужен |
| Конфиг | `os.Getenv` + своя структура `Config` | Проекту нужно ~8 переменных. Viper — из пушки по воробьям |
| Валидация (для админки) | [`go-playground/validator`](https://github.com/go-playground/validator) | Стандарт для проверки входящих JSON по тегам структуры |
| Авторизация админки | HTTP Basic Auth, пароль из переменной окружения | Один админ на всю систему — база пользователей, хеши и токены не нужны. `chi` даёт готовый `middleware.BasicAuth`, своего кода ноль. Подробности и оговорки — в разделе 4a |
| Тесты | стандартный `testing` + `net/http/httptest` | Ручки тестируются без запуска сервера |

Версия Go — актуальная стабильная (на момент старта — 1.26.x).

### Что сознательно не берём

- **Gin / Echo / Fiber.** Они вводят свой тип `*gin.Context` вместо стандартных
  `http.ResponseWriter` и `*http.Request`. Код становится непереносимым, а
  учишься ты фреймворку, а не языку. Для API из десяти ручек выгоды нет.
  (Если позже упрёшься в потолок — переход на Gin занимает вечер, потому что
  логика будет лежать в отдельных пакетах, а не в хендлерах.)
- **GORM** (самая популярная ORM). Прячет SQL за «магией», а когда запрос
  тормозит — непонятно, что она сгенерировала. С `sqlc` ты видишь ровно тот SQL,
  который уйдёт в базу, и попутно учишь SQL — он пригодится всегда.
- **JWT, bcrypt, таблица пользователей, OAuth.** Всё это нужно, когда
  пользователей много и они регистрируются сами. Здесь админ один, пароль
  задаётся на машине в переменной окружения — вся эта машинерия только добавила
  бы кода и мест, где можно ошибиться.

---

## 2. Структура проекта

Новый бэкенд кладём рядом со старым — `backend-go/`. Так заглушка продолжает
работать, пока Go-версия не догонит её по функциональности.

```
backend-go/
├── cmd/
│   └── api/
│       └── main.go            # точка входа: конфиг → БД → MinIO → роутер → сервер
├── internal/                  # "internal" = импортировать извне нельзя (правило Go)
│   ├── config/config.go       # чтение переменных окружения
│   ├── httpapi/
│   │   ├── router.go          # сборка маршрутов и middleware
│   │   ├── classes.go         # хендлеры занятий
│   │   ├── news.go            # хендлеры новостей
│   │   ├── ...                # services, staff, documents, gallery
│   │   ├── dto.go             # структуры ответов с json-тегами в camelCase
│   │   └── errors.go          # единый формат ошибки { "error": {...} }
│   ├── store/
│   │   ├── queries/*.sql      # SQL-запросы (вход для sqlc)
│   │   └── gen/               # сгенерированный sqlc код — не редактируем руками
│   └── media/
│       └── minio.go           # клиент MinIO: upload, delete, presigned URL
├── migrations/                # 0001_init.sql, 0002_... (goose)
├── seed/
│   └── main.go                # разовый импорт данных и фото из заглушки
├── sqlc.yaml
├── docker-compose.yml         # postgres + minio (+ сам api)
├── Makefile
└── .env.example
```

Почему так: `cmd/` — запускаемые программы, `internal/` — вся логика,
`httpapi` ничего не знает про SQL, `store` ничего не знает про HTTP. Это
позволяет тестировать по частям и не бояться менять одно, не сломав другое.

---

## 3. Схема базы данных

Главное решение: **фотографии живут в отдельной таблице `media`**, а сущности
ссылаются на неё. Причина видна уже в текущих данных — одно и то же фото
используется в нескольких местах (например `190.jpg` — у «Столярки»,
«Мастериц» и «Изолепки»). Если хранить путь строкой в каждой таблице, при
замене файла придётся править его в пяти местах.

```sql
-- Файл в MinIO. object_key — путь внутри бакета, например 'photos/090.jpg'
create table media (
  id          bigserial primary key,
  object_key  text not null unique,
  mime        text not null,             -- 'image/jpeg', 'application/pdf'
  size_bytes  bigint not null,
  width       int,                       -- для картинок, чтобы фронт мог зарезервировать место
  height      int,
  alt         text not null default '',  -- подпись для screen reader
  created_at  timestamptz not null default now()
);

-- Занятия
create table classes (
  id               bigserial primary key,
  title            text not null,
  age_range        text not null,
  description      text not null,
  schedule         text not null,
  price            text not null,        -- строкой: тарифы нерегулярные
  emoji            text not null,
  color            char(7) not null,     -- HEX: '#e46129'
  cover_media_id   bigint references media(id) on delete set null,  -- главное фото
  full_description text[] not null default '{}',   -- абзацы
  highlights       text[] not null default '{}',
  position         int not null default 0,         -- порядок вывода на сайте
  is_published     boolean not null default true,
  created_at       timestamptz not null default now(),
  updated_at       timestamptz not null default now()
);

-- Список фото занятия (галерея). Порядок — position
create table class_photos (
  class_id  bigint not null references classes(id) on delete cascade,
  media_id  bigint not null references media(id)   on delete cascade,
  position  int not null default 0,
  primary key (class_id, media_id)
);

-- Новости
create table news (
  id             bigserial primary key,
  title          text not null,
  published_on   date not null,          -- дата без времени, как в контракте
  category       text not null,
  excerpt        text not null,
  emoji          text not null,
  cover_media_id bigint references media(id) on delete set null,
  content        text[] not null default '{}',
  is_published   boolean not null default true,
  created_at     timestamptz not null default now(),
  updated_at     timestamptz not null default now()
);

create table news_photos (
  news_id  bigint not null references news(id)  on delete cascade,
  media_id bigint not null references media(id) on delete cascade,
  position int not null default 0,
  primary key (news_id, media_id)
);

-- Услуги
create table services (
  id             bigserial primary key,
  title          text not null,
  tagline        text not null,
  description    text not null,
  features       text[] not null default '{}',
  price          text not null,
  emoji          text not null,
  path           text not null,          -- '/special/birthday' — маршрут во фронте
  color          char(7) not null,
  cover_media_id bigint references media(id) on delete set null,
  detail         jsonb,                  -- блоки страницы услуги, см. контракт
  position       int not null default 0,
  is_published   boolean not null default true
);

create table service_photos (
  service_id bigint not null references services(id) on delete cascade,
  media_id   bigint not null references media(id)    on delete cascade,
  position   int not null default 0,
  primary key (service_id, media_id)
);

-- Сотрудники
create table staff (
  id           bigserial primary key,
  name         text not null,
  role         text not null,
  bio          text not null,
  initials     text not null,            -- запасной вариант, если фото нет
  color        char(7) not null,
  photo_media_id bigint references media(id) on delete set null,
  position     int not null default 0,
  is_published boolean not null default true
);

-- Документы. Сам PDF — тоже объект в MinIO
create table documents (
  id           bigserial primary key,
  title        text not null,
  description  text not null,
  file_media_id bigint not null references media(id) on delete restrict,
  position     int not null default 0,
  is_published boolean not null default true
);

-- Галерея
create table gallery (
  id       bigserial primary key,
  media_id bigint not null references media(id) on delete cascade,
  title    text not null default '',
  position int not null default 0
);

-- Отзывы (сейчас в frontend/src/data/reviews.ts)
create table reviews (
  id           bigserial primary key,
  author       text not null,
  role         text not null,
  text         text not null,
  rating       smallint not null check (rating between 1 and 5),
  position     int not null default 0,
  is_published boolean not null default true
);
```

Заметки по схеме:

- **`text[]` для абзацев** (`content`, `full_description`, `highlights`,
  `features`) — pgx маппит массив Postgres прямо в `[]string`. Отдельная таблица
  «абзацев» тут была бы оверинжинирингом.
- **`detail jsonb`** у услуг — структура блоков у каждой услуги разная
  (пакеты / форматы / смены), см. контракт. Загонять это в таблицы рано.
- **`position`** — числовой порядок сортировки. Сейчас порядок задан порядком
  элементов в питоновском списке; в БД порядок строк не гарантирован, нужна явная
  колонка. Сортировка: `order by position, id`.
- **`is_published`** — чтобы в админке можно было готовить новость и не
  показывать её на сайте. Публичные ручки всегда фильтруют по `is_published`.
- **`fileType` и `fileSize`** из контракта в БД **не храним** — считаем в Go из
  `media.mime` и `media.size_bytes` и форматируем в «PDF» и «1,4 МБ».
- `id` — `bigserial`, то есть целые числа. Контракт уже на числовых id, менять
  ничего не нужно.

---

## 4. Работа с MinIO

**Как устроено:** файл лежит в MinIO (бакет `media`, ключ вида
`photos/2026/07/uuid.jpg`), в Postgres — только `object_key`. Бакет **приватный**.
На каждый запрос API генерирует *presigned URL* — временную подписанную ссылку.

Это ровно то, что уже описано в контракте и заложено в заглушке
(функция `media_url()` в `backend-stub/app/main.py`), так что фронт снова не меняется.

```go
// Ссылка живёт 1 час. Генерация — локальная подпись HMAC,
// сетевого запроса в MinIO не делается, поэтому это дёшево
// даже когда в ответе 40 фото.
url, err := mc.PresignedGetObject(ctx, bucket, objectKey, time.Hour, nil)
```

Два подводных камня, которые лучше решить сразу:

1. **Кэш браузера.** Если ссылка каждый раз новая, браузер каждый раз качает
   картинку заново. Лечится округлением времени жизни: генерировать ссылку с
   истечением «в начале следующего часа», тогда в течение часа URL совпадает
   и кэш работает. Хелпер `media.URL()` инкапсулирует это.
2. **Приватный бакет vs публичный.** Presigned нужен для документов и
   персональных данных. Для публичных фото центра можно сделать бакет
   read-only-публичным и отдавать постоянные URL — быстрее и кэшируется идеально.
   Разумный компромисс: `photos/` — публичный префикс, `docs/` — приватный с
   presigned. Решение можно принять позже, контракт не меняется.

При загрузке через админку: проверяем MIME и размер, генерируем ключ через UUID
(имя файла от пользователя в ключ не подставляем), пережимаем в 2–3 размера
(например 480/960/1920 по ширине) — библиотека
[`disintegration/imaging`](https://github.com/disintegration/imaging) или
`h2non/bimg`. Тогда карточки не будут тянуть фото на 3 МБ.

---

## 4a. Авторизация админки

Админ один, пароль лежит в переменной окружения на сервере. Значит, ни таблицы
пользователей, ни хеширования, ни токенов не нужно.

### Вариант по умолчанию: HTTP Basic Auth

Самое короткое, что вообще возможно, — встроенный middleware `chi`:

```go
admin := chi.NewRouter()
admin.Use(middleware.BasicAuth("admin", map[string]string{
    cfg.AdminUser: cfg.AdminPassword, // оба значения из окружения
}))
admin.Post("/media", h.UploadMedia)
admin.Put("/classes/{id}", h.UpdateClass)
// ...
r.Mount("/api/admin", admin)
```

Всё: браузер сам покажет окно ввода логина и пароля и будет подставлять его в
последующие запросы. Своего кода — ноль, зависимостей — ноль.

Переменные окружения:

```
ADMIN_USER=admin
ADMIN_PASSWORD=<длинный случайный пароль>
```

### Что обязательно учесть

1. **Только HTTPS.** Basic Auth передаёт пароль в каждом запросе, закодированный
   в base64 — это **не шифрование**, а просто другая запись того же текста. По
   HTTP пароль уходит открытым. Поэтому в проде nginx/Caddy с сертификатом
   обязателен, и `/api/admin/*` стоит вообще не пускать без TLS.
2. **Пароль — длинный и случайный** (20+ символов из генератора). Он не
   хешируется и защищён только своей длиной. Никаких «Admin2026».
3. **`.env` не в гите.** Корневого `.gitignore` в проекте сейчас нет — при
   создании `backend-go/` сразу добавить туда `.env`. В `.env.example` кладём
   только имена переменных с пустыми значениями.
4. **Ограничение попыток.** Чтобы пароль нельзя было перебирать, добавить
   [`go-chi/httprate`](https://github.com/go-chi/httprate) на группу
   `/api/admin` — это одна строка: `admin.Use(httprate.LimitByIP(20, time.Minute))`.
5. **Сравнение в постоянном времени.** Если решишь писать проверку сам, а не
   брать готовый middleware, сравнивай не через `==`, а через
   `subtle.ConstantTimeCompare` — обычное сравнение строк завершается на первом
   несовпавшем символе, и по времени ответа пароль теоретически можно угадывать
   посимвольно.

### Когда понадобится нормальная форма логина

У Basic Auth два минуса: некрасивое системное окно браузера вместо своей формы
и невозможность «выйти» (браузер помнит пароль до закрытия). Если админка будет
на Vue и захочется свою страницу входа с кнопкой «Выйти» — апгрейд делается без
новых зависимостей:

- `POST /api/admin/login` — принимает пароль, сверяет с переменной окружения,
  в ответ ставит cookie: `HttpOnly`, `Secure`, `SameSite=Strict`, срок ~7 дней;
- значение cookie — строка `срок_истечения` плюс её HMAC-подпись секретом из
  окружения (`crypto/hmac` + `crypto/sha256`, обе — стандартная библиотека).
  Подпись не даёт подделать cookie, а срок внутри — использовать её вечно;
- middleware проверяет подпись и срок;
- `POST /api/admin/logout` — стирает cookie.

Это ~40 строк и одна дополнительная переменная окружения `SESSION_SECRET`.
Полноценный JWT здесь не нужен: JWT решает задачу «несколько сервисов доверяют
одному токену без обращения к базе», а у нас один сервис и один пользователь.

---

## 5. Маппинг ручек

Все ручки — только чтение, точно как сейчас. Формат ответа — байт в байт как в
[api-contracts.md](api-contracts.md), ключи в camelCase (задаются json-тегами в
`dto.go`).

| Метод | Путь | SQL | Заметки |
|-------|------|-----|---------|
| GET | `/api/news` | `news` + cover + `news_photos` | `order by published_on desc` |
| GET | `/api/news/{id}` | то же, одна запись | 404, если нет или не опубликована |
| GET | `/api/classes` | `classes` + cover + `class_photos` | `order by position, id` |
| GET | `/api/classes/{id}` | то же, одна запись | |
| GET | `/api/services` | `services` + фото | |
| GET | `/api/services/{id}` | то же + `detail` | |
| GET | `/api/staff` | `staff` + `photo_media_id` | |
| GET | `/api/documents` | `documents` + `media` | `fileType`/`fileSize` считаем в Go |
| GET | `/api/gallery` | `gallery` + `media` | |
| GET | `/api/reviews` | `reviews` | новая ручка, сейчас данные во фронте |
| GET | `/api/site` | `site_settings` (key/value или одна строка) | новая, опционально |
| GET | `/api/healthz` | — | для мониторинга и Docker healthcheck |

**Важная деталь — избегаем N+1.** Наивно: получить 14 занятий, потом на каждое
сделать отдельный запрос за фото — 15 запросов. Правильно: один запрос с
`array_agg`, который сразу собирает ключи фото в массив:

```sql
select c.*,
       cover.object_key as cover_key,
       coalesce(array_agg(m.object_key order by cp.position)
                filter (where m.id is not null), '{}') as photo_keys
from classes c
left join media cover on cover.id = c.cover_media_id
left join class_photos cp on cp.class_id = c.id
left join media m on m.id = cp.media_id
where c.is_published
group by c.id, cover.object_key
order by c.position, c.id;
```

Presigned-ссылки навешиваются уже в слое `httpapi`, когда данные превращаются в
DTO. `store` возвращает ключи объектов, а не URL — так БД-слой не зависит от MinIO.

---

## 6. Этапы работ

### Этап 0. Каркас и инфраструктура

- [x] `docker-compose.yml`: `postgres:17`, `minio/minio` (порты 9000/9001 —
      API и веб-консоль), тома для данных
- [x] `.gitignore` в `backend-go/` с `.env` в первой строке
- [x] `.env.example` (только имена, без значений): `DATABASE_URL`,
      `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_BUCKET`,
      `MEDIA_PUBLIC_BASE_URL`, `PORT`, `CORS_ORIGINS`, `ADMIN_USER`,
      `ADMIN_PASSWORD`
- [x] Конфиг падает на старте с понятной ошибкой, если обязательной переменной
      нет. Пустой `ADMIN_PASSWORD` — это открытая админка, такое лучше поймать
      сразу, а не в проде
- [x] `go mod init`, подключить chi, pgx, minio-go, goose
- [x] `main.go`: чтение конфига, пул pgx, клиент MinIO, chi-роутер с
      middleware (`RequestID`, `RealIP`, `Recoverer`, `Timeout` из
      `chi/middleware`, CORS из отдельного пакета `go-chi/cors`, свой
      middleware логирования на slog), ручка `/api/healthz`,
      graceful shutdown по `SIGTERM`
- [x] `Makefile`: `make up`, `make migrate`, `make sqlc`, `make run`, `make test`

**Готово, когда:** `curl localhost:8080/api/healthz` отвечает `{"ok":true}`.

### Этап 1. Схема и миграции

- [x] `migrations/0001_init.sql` — таблицы из раздела 3
- [ ] `sqlc.yaml` + первые запросы в `internal/store/queries/`
- [ ] `make sqlc` генерирует код, он компилируется

**Готово, когда:** `goose up` создаёт схему на чистой базе, `goose down` её
откатывает.

### Этап 2. Импорт текущих данных

Самый приятный этап: контент уже есть и проверен, надо просто перенести.

- [ ] Снять дампы с работающей заглушки:
      `curl -s localhost:3001/api/classes > seed/classes.json` (и остальные)
- [ ] `seed/main.go`: залить `backend-stub/media/photos/*.jpg` в MinIO (заодно
      посчитать `width`/`height`/`size`), затем вставить записи, подставляя
      `media.id` вместо путей
- [ ] Сверить: количество записей и набор фото совпадают с заглушкой

**Готово, когда:** ответы Go-бэкенда и Python-заглушки совпадают с точностью до
домена в ссылках на фото. Проверить можно так:

```bash
diff <(curl -s localhost:3001/api/classes | jq -S 'map(del(.photo,.photos))') \
     <(curl -s localhost:8080/api/classes | jq -S 'map(del(.photo,.photos))')
```

### Этап 3. Публичные ручки

- [x] Все публичные ручки из раздела 5, включая 404 на несуществующий id и 400 на
      нечисловой
- [ ] `dto.go` с json-тегами в camelCase; необязательные поля — `omitempty`,
      чтобы отсутствующее фото не приезжало как `"photo": ""`
- [x] Единый формат ошибки: `{"error":{"code":"not_found","message":"..."}}`
- [ ] Тесты хендлеров на `httptest` + `golden`-файлы с ожидаемым JSON
- [x] Переключить прокси Vite (`frontend/vite.config.ts`) на порт Go-бэкенда

**Готово, когда:** сайт целиком работает на Go-бэкенде, во фронте не изменено
ни строки, кроме порта в конфиге Vite.

### Этап 4. Запись и админка

До этого этапа контент правится SQL-ом — это нормально для старта.

- [ ] Группа `/api/admin/*` под `middleware.BasicAuth` с паролем из окружения
      (см. раздел 4a) + `httprate` на перебор
- [ ] `POST /api/admin/media` — загрузка файла (multipart), проверка MIME и
      размера, генерация превью, запись в `media`
- [ ] CRUD на занятия / новости / услуги / сотрудников / документы / галерею /
      отзывы + `PATCH .../reorder` для порядка (`position`)
- [ ] Валидация входящих JSON через `validator`
- [ ] Только после этого — UI админки (отдельная задача, скорее всего второй
      Vue-роут `/admin` или отдельное приложение)

### Этап 5. Прод

- [ ] Multi-stage `Dockerfile` (сборка `CGO_ENABLED=0` → `FROM gcr.io/distroless/static`,
      образ выходит ~15 МБ)
- [ ] Миграции накатываются при старте контейнера
- [ ] nginx/Caddy впереди: HTTPS, отдача статики фронта, проксирование `/api`
- [ ] Бэкапы: `pg_dump` по cron + `mc mirror` бакета
- [ ] Логи в JSON (`slog.NewJSONHandler`), метрики Prometheus — по желанию
- [ ] Удалить `backend-stub/` (Python-заглушку) и обновить README

---

## 7. Что не менять

Эти вещи уже приведены к финальному виду, и Go-бэкенд должен их сохранить:

- `id` — целые числа;
- `color` — HEX-строка (фронт сам считает оттенки через `color-mix`);
- `photo` — главное фото, `photos` — список; ссылки **абсолютные**;
- ключи JSON — camelCase;
- списки — плоские массивы, без обёртки `{items, total}` (пагинация не нужна,
  данных мало).

---

## 8. Порядок изучения (чтобы не утонуть)

Раз опыта в Go пока мало, полезно идти в таком порядке — каждый шаг даёт
работающий результат:

1. **Tour of Go** — только базовый синтаксис, структуры, интерфейсы, `error`.
   Горутины и каналы для этого проекта почти не нужны.
2. Ручка `/api/healthz` на чистом `net/http`, без БД. Понять `http.Handler`,
   `json.NewEncoder(w).Encode(v)`, json-теги структур.
3. Подключить pgx и сделать одну ручку `/api/gallery` (самая простая: одна
   таблица, никаких массивов).
4. Добавить MinIO и presigned-ссылки в ту же ручку.
5. Дальше остальные ручки по образцу — они однотипные.

Ключевая идея, из-за которой Go поначалу кажется многословным: ошибки не
бросаются, а возвращаются, и их проверяют на каждом шаге
(`if err != nil { return err }`). Это не бойлерплейт «просто так» — это то, что
делает Go-сервисы предсказуемыми.
