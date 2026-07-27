"""API-заглушка сайта «Мы Друзья» на Litestar.

Пока это «фейковый» бэкенд: данные лежат рядом в app/data/*.py (те же объекты,
что раньше были захардкожены во фронте) и отдаются как JSON. Когда появится
настоящая база, достаточно заменить источники данных ниже на запросы к БД —
контракты ручек не изменятся.

Контракты ручек описаны в docs/api-contracts.md.
"""

from litestar import Litestar, get
from litestar.config.cors import CORSConfig
from litestar.exceptions import NotFoundException

from app.data.classes import classes
from app.data.documents import documents
from app.data.gallery import gallery
from app.data.news import news
from app.data.services import special_offers
from app.data.staff import staff

Item = dict[str, object]


def _find(items: list[Item], item_id: str, not_found: str) -> Item:
    """Найти элемент по id или отдать 404."""
    for item in items:
        if item["id"] == item_id:
            return item
    raise NotFoundException(detail=not_found)


# --- Новости ---------------------------------------------------------------
@get("/api/news")
async def list_news() -> list[Item]:
    return news


@get("/api/news/{item_id:str}")
async def get_news(item_id: str) -> Item:
    return _find(news, item_id, "Новость не найдена")


# --- Занятия ---------------------------------------------------------------
@get("/api/classes")
async def list_classes() -> list[Item]:
    return classes


@get("/api/classes/{item_id:str}")
async def get_class(item_id: str) -> Item:
    return _find(classes, item_id, "Занятие не найдено")


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
    return gallery


# --- Проверка живости: GET /api --------------------------------------------
@get("/api")
async def health() -> dict[str, object]:
    return {"ok": True, "service": "children-center-backend"}


# В дев-режиме фронт ходит через прокси Vite, но с CORS ручки работают и при
# прямом обращении, и из прод-сборки, размещённой на другом домене.
cors_config = CORSConfig(allow_origins=["*"])

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
    ],
    cors_config=cors_config,
)
