/* Типы данных сайта.
   Когда появится бэкенд, ответы API должны совпадать с этими типами —
   и фронтенд заработает без переделок. */

export interface NavLink {
  label: string
  to: string
  /** Вложенные ссылки (для выпадающих меню) */
  children?: NavLink[]
}

/** Контент главной страницы. */
export interface HomeContent {
  /** Фото, которые по очереди показываются в главном баннере. */
  heroPhotos: string[]
}

/** Акцентный цвет карточки — HEX-строка с бэкенда, например '#1a6757'.
    Светлую подложку и затемнённый оттенок фронт считает сам (см. useAccent). */
export type HexColor = string

export interface ClassItem {
  id: number
  title: string
  ageRange: string
  description: string
  schedule: string
  price: string
  emoji: string
  /** Главное фото занятия. Пока не задано — показывается заглушка. */
  photo?: string
  /** Дополнительные фото с занятий — мини-галерея на странице занятия. */
  photos?: string[]
  color: HexColor
  /** Развёрнутое описание для отдельной страницы занятия (абзацы). */
  fullDescription?: string[]
  /** Что получит ребёнок — список на странице занятия. */
  highlights?: string[]
}

export interface StaffMember {
  id: number
  name: string
  role: string
  bio: string
  /** Инициалы для заглушки-аватара, если фото нет */
  initials: string
  color: HexColor
  /** Фото сотрудника — показывается в кружке аватара вместо инициалов. */
  photo?: string
}

export interface Review {
  id: number
  author: string
  role: string
  text: string
  rating: number
}

export interface NewsPost {
  id: number
  title: string
  date: string // ISO-формат: YYYY-MM-DD
  category: string
  excerpt: string
  emoji: string
  /** Обложка новости. Не задано — показывается заглушка. */
  photo?: string
  /** Остальные фото к новости — галерея под текстом. */
  photos?: string[]
  /** Полный текст новости (абзацы) для отдельной страницы. */
  content?: string[]
}

export interface DocumentItem {
  id: number
  title: string
  description: string
  fileType: string // например 'PDF'
  fileSize: string // например '1,2 МБ'
  url: string
}

export interface GalleryItem {
  id: number
  /** Абсолютная ссылка на фото с бэкенда */
  image: string
}

export interface SpecialOffer {
  id: number
  title: string
  tagline: string
  description: string
  features: string[]
  price: string
  emoji: string
  to: string
  color: HexColor
  /** Главное фото услуги */
  photo?: string
  /** Галерея к услуге */
  photos?: string[]
}
