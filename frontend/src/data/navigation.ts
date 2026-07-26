import type { NavLink } from './types'

/* Главное меню. Один источник для шапки, мобильного меню и подвала. */
export const mainNav: NavLink[] = [
  { label: 'Главная', to: '/' },
  {
    label: 'О центре',
    to: '/about',
    children: [
      { label: 'О нас', to: '/about' },
      { label: 'Галерея', to: '/about/gallery' },
      { label: 'Сотрудники', to: '/about/staff' },
      { label: 'Отзывы', to: '/about/reviews' },
      { label: 'Документы', to: '/about/documents' },
    ],
  },
  { label: 'Занятия', to: '/classes' },
  { label: 'Новости', to: '/news' },
  {
    label: 'Услуги',
    to: '/special/birthday',
    children: [
      { label: 'День рождения', to: '/special/birthday' },
      { label: 'Мастер-классы', to: '/special/masterclasses' },
      { label: 'Летние каникулы', to: '/special/summer' },
    ],
  },
  { label: 'Контакты', to: '/contacts' },
]

/* Под-навигация раздела «О центре» */
export const aboutNav: NavLink[] = [
  { label: 'О нас', to: '/about' },
  { label: 'Галерея', to: '/about/gallery' },
  { label: 'Сотрудники', to: '/about/staff' },
  { label: 'Отзывы', to: '/about/reviews' },
  { label: 'Документы', to: '/about/documents' },
]
