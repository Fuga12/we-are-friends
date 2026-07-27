"""API-заглушка сайта «Мы Друзья» на Litestar.

Пока это «фейковый» бэкенд: данные лежат рядом в app/data/*.py (те же объекты,
что раньше были захардкожены во фронте) и отдаются как JSON. Когда появится
настоящая база, достаточно заменить источники данных ниже на запросы к БД —
контракты ручек не изменятся.

Контракты ручек описаны в docs/api-contracts.md.
"""

import os
from pathlib import Path

from litestar import Litestar, get
from litestar.config.cors import CORSConfig
from litestar.exceptions import NotFoundException
from litestar.static_files import create_static_files_router

from app.data.classes import classes
from app.data.documents import documents
from app.data.gallery import gallery
from app.data.news import news
from app.data.services import special_offers
from app.data.staff import staff

Item = dict[str, object]

# Папка с медиафайлами (фото центра). Путь считаем от этого файла, чтобы
# сервер работал из любой директории: backend/media.
MEDIA_DIR = Path(__file__).resolve().parent.parent / "media"

# Базовый адрес, из которого собираются АБСОЛЮТНЫЕ ссылки на фото.
# В деве это сам бэкенд; в проде задаётся переменной окружения MEDIA_BASE_URL
# (например, адрес S3-бакета или CDN).
MEDIA_BASE_URL = os.environ.get("MEDIA_BASE_URL", "http://localhost:3001").rstrip("/")


def media_url(path: str) -> str:
    """Собрать абсолютную ссылку на файл.

    Заглушка просто приклеивает базовый адрес: /media/photos/090.jpg ->
    http://localhost:3001/media/photos/090.jpg.

    В проде здесь будет генерация ПОДПИСАННОЙ ссылки на объект в S3
    (presigned URL с истечением срока) — контракт от этого не меняется:
    в JSON так же приходит готовая абсолютная ссылка, которую фронт вставляет.
    """
    return f"{MEDIA_BASE_URL}{path}"


def _absolutize(item: Item, *fields: str) -> Item:
    """Вернуть копию элемента, где медиа-поля превращены в абсолютные ссылки."""
    out = dict(item)
    for field in fields:
        value = out.get(field)
        if isinstance(value, str) and value.startswith("/media/"):
            out[field] = media_url(value)
    return out


def _find(items: list[Item], item_id: str, not_found: str) -> Item:
    """Найти элемент по id или отдать 404."""
    for item in items:
        if item["id"] == item_id:
            return item
    raise NotFoundException(detail=not_found)


# --- Новости ---------------------------------------------------------------
@get("/api/news")
async def list_news() -> list[Item]:
    return [_absolutize(n, "photo") for n in news]


@get("/api/news/{item_id:str}")
async def get_news(item_id: str) -> Item:
    return _absolutize(_find(news, item_id, "Новость не найдена"), "photo")


# --- Занятия ---------------------------------------------------------------
@get("/api/classes")
async def list_classes() -> list[Item]:
    return [_absolutize(c, "photo") for c in classes]


@get("/api/classes/{item_id:str}")
async def get_class(item_id: str) -> Item:
    return _absolutize(_find(classes, item_id, "Занятие не найдено"), "photo")


# --- Услуги ----------------------------------------------------------------
@get("/api/services")
async def list_services() -> list[Item]:
    return special_offers


@get("/api/services/{item_id:str}")
async def get_service(item_id: str) -> Item:
    return _find(special_offers, item_id, "Услуга не найдена")


# --- Списки без отдельных страниц ------------------------------------------
@get("/api/staff")
async def list_staff() -> list[Item]:
    return staff


@get("/api/documents")
async def list_documents() -> list[Item]:
    return documents


@get("/api/gallery")
async def list_gallery() -> list[Item]:
    return [_absolutize(g, "image") for g in gallery]


# --- Проверка живости: GET /api --------------------------------------------
@get("/api")
async def health() -> dict[str, object]:
    return {"ok": True, "service": "children-center-backend"}


# В дев-режиме фронт ходит через прокси Vite, но с CORS ручки работают и при
# прямом обращении, и из прод-сборки, размещённой на другом домене.
cors_config = CORSConfig(allow_origins=["*"])

# Раздача медиафайлов: GET /media/photos/090.jpg -> backend/media/photos/090.jpg
media_router = create_static_files_router(path="/media", directories=[MEDIA_DIR])

app = Litestar(
    route_handlers=[
        list_news,
        get_news,
        list_classes,
        get_class,
        list_services,
        get_service,
        list_staff,
        list_documents,
        list_gallery,
        health,
        media_router,
    ],
    cors_config=cors_config,
)
