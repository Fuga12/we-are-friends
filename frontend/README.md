# Frontend

Vue 3, TypeScript, Vite, Vue Router и PrimeVue.

## Проверка и сборка

Нужен Node.js 22 (та же версия используется в Dockerfile).

```bash
npm ci
npm run type-check
npm run build
```

`build` включает проверку типов и создаёт `dist/`.

## Запуск сайта

Для работы с Go API и MinIO используйте Docker Compose из корня репозитория — [инструкция](../README.md). Nginx направляет `/api` в backend и `/media` в объектное хранилище.

Для работы только над интерфейсом можно запустить `npx vite --host 127.0.0.1`. Этот режим сам по себе не запускает backend и не настраивает прокси API.

В публичной версии фотографии заменены демонстрационной SVG-иллюстрацией.
