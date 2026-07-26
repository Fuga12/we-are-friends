/* Типы данных сайта.
   Когда появится бэкенд, ответы API должны совпадать с этими типами —
   и фронтенд заработает без переделок. */

export interface NavLink {
  label: string
  to: string
  /** Вложенные ссылки (для выпадающих меню) */
  children?: NavLink[]
}

export interface ClassItem {
  id: string
  title: string
  ageRange: string
  description: string
  schedule: string
  price: string
  emoji: string
  /** Ссылка на фото занятия. Пока не задано — показывается заглушка. */
  photo?: string
  color: 'teal' | 'orange' | 'amber' | 'red'
  /** Развёрнутое описание для отдельной страницы занятия (абзацы). */
  fullDescription?: string[]
  /** Что получит ребёнок — список на странице занятия. */
  highlights?: string[]
}

export interface StaffMember {
  id: string
  name: string
  role: string
  bio: string
  /** Инициалы для заглушки-аватара, пока нет фото */
  initials: string
  color: 'teal' | 'orange' | 'amber' | 'red'
}

export interface Review {
  id: string
  author: string
  role: string
  text: string
  rating: number
}

export interface NewsPost {
  id: string
  title: string
  date: string // ISO-формат: YYYY-MM-DD
  category: string
  excerpt: string
  emoji: string
  /** Обложка новости. Не задано — показывается заглушка. */
  photo?: string
  /** Полный текст новости (абзацы) для отдельной страницы. */
  content?: string[]
}

export interface DocumentItem {
  id: string
  title: string
  description: string
  fileType: string // например 'PDF'
  fileSize: string // например '1,2 МБ'
  url: string
}

export interface GalleryItem {
  id: string
  /** Путь к фото в public/photos */
  image: string
  /** Подпись (необязательно) */
  title?: string
}

export interface SpecialOffer {
  id: string
  title: string
  tagline: string
  description: string
  features: string[]
  price: string
  emoji: string
  to: string
  color: 'teal' | 'orange' | 'amber' | 'red'
}
