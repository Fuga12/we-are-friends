/* Слой работы с бэкендом.
   Один источник всех запросов к API — компоненты не знают про fetch и URL,
   а зовут готовые функции вроде getNews(). Типы ответов — те же, что в
   src/data/types.ts, поэтому шаблоны компонентов не меняются.

   Запросы идут на относительный /api — в разработке Vite проксирует их на
   бэкенд (см. vite.config.ts). */
import type {
  ClassItem,
  DocumentItem,
  GalleryItem,
  NewsPost,
  SpecialOffer,
  StaffMember,
} from '@/data/types'

const BASE = '/api'

/** Базовый GET-запрос: сходить по адресу и вернуть распарсенный JSON.
    Если сервер ответил ошибкой (например 404) — бросаем исключение. */
async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE}${path}`)
  if (!res.ok) {
    throw new Error(`Запрос ${path} завершился со статусом ${res.status}`)
  }
  return (await res.json()) as T
}

/* id везде числовой. Из адреса страницы он приходит строкой
   (route.params.id), поэтому функции принимают и то и другое — в URL запроса
   всё равно подставляется то же самое число. */
type Id = number | string

// Новости
export const getNews = () => get<NewsPost[]>('/news')
export const getNewsPost = (id: Id) => get<NewsPost>(`/news/${id}`)

// Занятия
export const getClasses = () => get<ClassItem[]>('/classes')
export const getClass = (id: Id) => get<ClassItem>(`/classes/${id}`)

// Услуги
export const getServices = () => get<SpecialOffer[]>('/services')
export const getService = (id: Id) => get<SpecialOffer>(`/services/${id}`)

// Прочие списки
export const getStaff = () => get<StaffMember[]>('/staff')
export const getDocuments = () => get<DocumentItem[]>('/documents')
export const getGallery = () => get<GalleryItem[]>('/gallery')
