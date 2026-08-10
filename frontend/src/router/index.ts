import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

/*
  Страницы импортируем обычным образом (не через () => import()).
  Сайт небольшой, весь код грузится сразу одним лёгким бандлом — зато
  переходы между страницами мгновенные, без подгрузки и «мигания».

  meta.title — заголовок вкладки браузера, ставим его в afterEach ниже.
  Раздел «О центре» сделан вложенным (nested): у него есть свой под-навигатор
  и общий каркас AboutLayout с <router-view> внутри.
*/
import HomeView from '@/views/HomeView.vue'
import AboutLayout from '@/views/about/AboutLayout.vue'
import AboutView from '@/views/about/AboutView.vue'
import GalleryView from '@/views/about/GalleryView.vue'
import StaffView from '@/views/about/StaffView.vue'
import ReviewsView from '@/views/about/ReviewsView.vue'
import DocumentsView from '@/views/about/DocumentsView.vue'
import NewsView from '@/views/NewsView.vue'
import NewsDetailView from '@/views/NewsDetailView.vue'
import ClassesView from '@/views/ClassesView.vue'
import ClassDetailView from '@/views/ClassDetailView.vue'
import ContactsView from '@/views/ContactsView.vue'
import BirthdayView from '@/views/special/BirthdayView.vue'
import MasterclassesView from '@/views/special/MasterclassesView.vue'
import SummerView from '@/views/special/SummerView.vue'
import NotFoundView from '@/views/NotFoundView.vue'
import AdminView from '@/views/admin/AdminView.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
    meta: { title: 'Мы друзья — детский центр в Петрозаводске' },
  },
  { path: '/admin', name: 'admin', component: AdminView, meta: { title: 'Админка', admin: true } },

  // Раздел «О центре» с вложенными подстраницами
  {
    path: '/about',
    component: AboutLayout,
    children: [
      { path: '', name: 'about', component: AboutView, meta: { title: 'О центре' } },
      { path: 'gallery', name: 'gallery', component: GalleryView, meta: { title: 'Галерея' } },
      { path: 'staff', name: 'staff', component: StaffView, meta: { title: 'Сотрудники' } },
      { path: 'reviews', name: 'reviews', component: ReviewsView, meta: { title: 'Отзывы' } },
      {
        path: 'documents',
        name: 'documents',
        component: DocumentsView,
        meta: { title: 'Документы' },
      },
    ],
  },

  { path: '/news', name: 'news', component: NewsView, meta: { title: 'Новости' } },
  // Отдельная страница новости. :id — динамический параметр из адреса.
  { path: '/news/:id', name: 'news-detail', component: NewsDetailView },
  { path: '/classes', name: 'classes', component: ClassesView, meta: { title: 'Занятия' } },
  // Отдельная страница занятия.
  { path: '/classes/:id', name: 'class-detail', component: ClassDetailView },
  { path: '/contacts', name: 'contacts', component: ContactsView, meta: { title: 'Контакты' } },

  // Специальные услуги
  {
    path: '/special/birthday',
    name: 'birthday',
    component: BirthdayView,
    meta: { title: 'День рождения ребёнка' },
  },
  {
    path: '/special/masterclasses',
    name: 'masterclasses',
    component: MasterclassesView,
    meta: { title: 'Мастер-классы' },
  },
  {
    path: '/special/summer',
    name: 'summer',
    component: SummerView,
    meta: { title: 'Летние каникулы' },
  },
  { path: '/special/summer/:id', name: 'summer-shift-detail', component: ClassDetailView, meta: { title: 'Летняя смена' } },

  // Страница 404 — ловит все несуществующие адреса
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundView,
    meta: { title: 'Страница не найдена' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  // При переходе на новую страницу — прокручиваем наверх.
  // При переходе назад/вперёд — возвращаем прежнюю позицию.
  // Если есть якорь (#id) — прокручиваем к нему.
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, behavior: 'smooth', top: 90 }
    return { top: 0 }
  },
})

// Обновляем заголовок вкладки браузера под каждую страницу
const SITE_NAME = 'Мы друзья'
router.afterEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? (to.name === 'home' ? title : `${title} · ${SITE_NAME}`) : SITE_NAME
})

export default router
