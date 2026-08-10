<script setup lang="ts">
import BaseButton from '@/components/ui/BaseButton.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import PhotoFrame from '@/components/ui/PhotoFrame.vue'
import ClassCard from '@/components/cards/ClassCard.vue'
import SpecialOfferCard from '@/components/cards/SpecialOfferCard.vue'
import ReviewCard from '@/components/cards/ReviewCard.vue'
import NewsCard from '@/components/cards/NewsCard.vue'
import CtaSection from '@/components/sections/CtaSection.vue'

import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { site } from '@/data/site'
import { reviews } from '@/data/reviews'
import { specialOffers } from '@/data/services'
import { useAsync } from '@/composables/useAsync'
import { getClasses, getHome, getNews } from '@/api'

// Данные-витрину тянем с бэкенда. Пока грузятся — списки пустые,
// секции просто дорисуются, как только придут данные.
const { data: classesData } = useAsync(getClasses)
const { data: homeData } = useAsync(getHome)
const { data: newsData } = useAsync(getNews)

// На главной показываем только часть данных — «витрину»
const featuredClasses = computed(() => (classesData.value ?? []).slice(0, 3))
const latestNews = computed(() => (newsData.value ?? []).slice(0, 3))
const featuredReviews = reviews.slice(0, 3)

const heroPhotoIndex = ref(0)
const heroPhotos = computed(() => homeData.value?.heroPhotos ?? [])
const readyHeroPhotos = ref<string[]>([])
let heroPhotoTimer: ReturnType<typeof setInterval> | undefined
let heroPhotoPreloadRun = 0

function stopHeroPhotoTimer() {
  if (heroPhotoTimer) {
    clearInterval(heroPhotoTimer)
    heroPhotoTimer = undefined
  }
}

function startHeroPhotoTimer() {
  stopHeroPhotoTimer()
  if (readyHeroPhotos.value.length < 2) return
  heroPhotoTimer = setInterval(() => {
    heroPhotoIndex.value = (heroPhotoIndex.value + 1) % readyHeroPhotos.value.length
  }, 7000)
}

function preloadHeroPhoto(src: string): Promise<string | undefined> {
  return new Promise((resolve) => {
    const image = new Image()
    image.onload = () => resolve(src)
    image.onerror = () => resolve(undefined)
    image.src = src
  })
}

watch(heroPhotos, async (photos) => {
  const preloadRun = ++heroPhotoPreloadRun
  stopHeroPhotoTimer()
  readyHeroPhotos.value = []
  heroPhotoIndex.value = 0

  const loadedPhotos = await Promise.all(photos.map(preloadHeroPhoto))
  if (preloadRun !== heroPhotoPreloadRun) return

  readyHeroPhotos.value = loadedPhotos.filter((src): src is string => src !== undefined)
  startHeroPhotoTimer()
})

onBeforeUnmount(stopHeroPhotoTimer)

const values = [
  { emoji: '🤗', title: 'Как семья', text: 'Тёплое и уважительное отношение к каждому ребёнку.' },
  { emoji: '🛠️', title: 'Настоящие навыки', text: 'Работаем настоящими инструментами и материалами, а не игрушками.' },
  { emoji: '🎨', title: 'Дизайн-центр', text: 'Художественное конструирование: делаем реальные вещи своими руками.' },
  { emoji: '🏆', title: 'Результат', text: 'Участвуем и побеждаем в конкурсах, фестивалях и модных показах.' },
]

const stats = [
  { value: '15', label: 'направлений' },
  { value: '7', label: 'лет работы' },
  { value: '5,0★', label: 'рейтинг в Яндекс Картах' },
  { value: '2', label: 'площадки в городе' },
]
</script>

<template>
  <!-- ХЕДЕР-БАННЕР -->
  <section class="hero">
    <div class="container hero__inner">
      <div class="hero__content">
        <span class="eyebrow">{{ site.tagline }}</span>
        <h1 class="hero__title">
          Место, где детям <span class="hero__accent">хорошо</span> расти
        </h1>
        <p class="hero__lead">{{ site.description }}</p>
        <div class="hero__actions">
          <BaseButton to="/classes" variant="primary" size="lg">Наши занятия</BaseButton>
          <BaseButton to="/contacts" variant="outline" size="lg">Записаться</BaseButton>
        </div>
        <div class="hero__stats">
          <div v-for="stat in stats" :key="stat.label" class="hero__stat">
            <span class="hero__stat-value">{{ stat.value }}</span>
            <span class="hero__stat-label">{{ stat.label }}</span>
          </div>
        </div>
      </div>

      <!-- Фото центра с парой декоративных эмодзи-бейджей -->
      <div class="hero__art">
        <div class="blob blob--teal" aria-hidden="true"></div>
        <div class="blob blob--amber" aria-hidden="true"></div>
        <div class="hero__photo">
          <template v-if="readyHeroPhotos.length">
            <img
              v-for="(src, index) in readyHeroPhotos"
              :key="src"
              :src="src"
              alt="Дети в центре «Мы Друзья»"
              :class="{ 'hero__photo-image--active': index === heroPhotoIndex }"
            />
          </template>
          <PhotoFrame v-else ratio="4 / 5" emoji="📸" />
        </div>
        <div class="hero__badge hero__badge--1" aria-hidden="true">🎨</div>
        <div class="hero__badge hero__badge--2" aria-hidden="true">⭐</div>
      </div>
    </div>
  </section>

  <!-- ЦЕННОСТИ -->
  <section class="section">
    <div class="container">
      <SectionHeader
        eyebrow="Почему мы"
        title="Тёплый дом для развития"
        subtitle="Мы создаём атмосферу, в которой ребёнку хочется учиться, дружить и пробовать новое."
        center
      />
      <div class="grid" style="--min: 240px">
        <article v-for="value in values" :key="value.title" class="value-card">
          <span class="value-card__emoji" aria-hidden="true">{{ value.emoji }}</span>
          <h3 class="value-card__title">{{ value.title }}</h3>
          <p class="value-card__text">{{ value.text }}</p>
        </article>
      </div>
    </div>
  </section>

  <!-- ЗАНЯТИЯ -->
  <section class="section section--muted">
    <div class="container">
      <div class="section-head-row">
        <SectionHeader eyebrow="Занятия" title="Выберите своё направление" />
        <BaseButton to="/classes" variant="ghost">Все направления →</BaseButton>
      </div>
      <div class="grid" style="--min: 280px">
        <ClassCard v-for="item in featuredClasses" :key="item.id" :item="item" />
      </div>
    </div>
  </section>

  <!-- СПЕЦИАЛЬНЫЕ УСЛУГИ -->
  <section class="section">
    <div class="container">
      <SectionHeader
        eyebrow="Особенное"
        title="Праздники и каникулы"
        subtitle="Больше, чем занятия: помогаем сделать детство ярким и запоминающимся."
        center
      />
      <div class="grid" style="--min: 300px">
        <SpecialOfferCard v-for="offer in specialOffers" :key="offer.id" :offer="offer" />
      </div>
    </div>
  </section>

  <!-- ОТЗЫВЫ -->
  <section class="section section--soft-teal">
    <div class="container">
      <div class="section-head-row">
        <SectionHeader eyebrow="Отзывы" title="Что говорят родители" />
        <BaseButton to="/about/reviews" variant="ghost">Все отзывы →</BaseButton>
      </div>
      <div class="grid" style="--min: 300px">
        <ReviewCard v-for="review in featuredReviews" :key="review.id" :review="review" />
      </div>
    </div>
  </section>

  <!-- НОВОСТИ -->
  <section class="section">
    <div class="container">
      <div class="section-head-row">
        <SectionHeader eyebrow="Новости" title="Что у нас происходит" />
        <BaseButton to="/news" variant="ghost">Все новости →</BaseButton>
      </div>
      <div class="stack-lg">
        <NewsCard v-for="post in latestNews" :key="post.id" :post="post" />
      </div>
    </div>
  </section>

  <CtaSection />
</template>

<style scoped>
/* ---------- HERO ---------- */
.hero {
  position: relative;
  overflow: hidden;
  background: radial-gradient(
    120% 120% at 100% 0%,
    var(--color-highlight-soft) 0%,
    var(--color-bg) 45%
  );
}
.hero__inner {
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: var(--space-7);
  align-items: center;
  padding-block: clamp(2.5rem, 1.5rem + 5vw, 5.5rem);
}
.hero__title {
  margin-top: var(--space-3);
}
.hero__accent {
  color: var(--color-accent);
  position: relative;
  white-space: nowrap;
}
.hero__accent::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0.05em;
  height: 0.28em;
  background-color: var(--brand-amber);
  opacity: 0.4;
  border-radius: var(--radius-pill);
  z-index: -1;
}
.hero__lead {
  margin-top: var(--space-4);
  max-width: 46ch;
  font-size: var(--text-lg);
  color: var(--color-text-muted);
}
.hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
  margin-top: var(--space-6);
}
.hero__stats {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-6);
  margin-top: var(--space-7);
  padding-top: var(--space-5);
  border-top: 1px dashed var(--color-border);
}
.hero__stat {
  display: flex;
  flex-direction: column;
}
.hero__stat-value {
  font-family: var(--font-heading);
  font-size: var(--text-2xl);
  font-weight: var(--weight-bold);
  color: var(--color-primary);
}
.hero__stat-label {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}

/* Декоративная часть справа */
.hero__art {
  position: relative;
  aspect-ratio: 1;
  min-height: 320px;
}
.blob {
  position: absolute;
  border-radius: 42% 58% 63% 37% / 41% 44% 56% 59%;
  filter: blur(2px);
}
.blob--teal {
  inset: 8% 12% 20% 8%;
  background-color: var(--brand-teal);
  opacity: 0.12;
}
.blob--amber {
  inset: 22% 6% 8% 26%;
  background-color: var(--brand-amber);
  opacity: 0.18;
  border-radius: 58% 42% 37% 63% / 44% 56% 44% 56%;
}
/* Рамка с фото по центру декоративной области */
.hero__photo {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: min(72%, 320px);
  padding: 10px;
  background-color: var(--color-surface);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-lg);
  rotate: -2deg;
  aspect-ratio: 4 / 5;
  overflow: hidden;
}
.hero__photo :deep(.photo) {
  height: 100%;
  border-radius: var(--radius-lg);
}
.hero__photo img {
  position: absolute;
  inset: 10px;
  width: calc(100% - 20px);
  height: calc(100% - 20px);
  object-fit: cover;
  border-radius: var(--radius-lg);
  opacity: 0;
  transition: opacity 0.8s var(--ease);
}
.hero__photo img.hero__photo-image--active {
  opacity: 1;
}
/* Плавающие эмодзи-бейджи */
.hero__badge {
  position: absolute;
  display: grid;
  place-items: center;
  width: clamp(56px, 7vw, 76px);
  height: clamp(56px, 7vw, 76px);
  font-size: clamp(26px, 3.5vw, 38px);
  background-color: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  animation: float 6s var(--ease) infinite;
}
.hero__badge--1 {
  top: 6%;
  left: 4%;
  animation-delay: 0s;
}
.hero__badge--2 {
  bottom: 8%;
  right: 4%;
  animation-delay: 1.5s;
}
@keyframes float {
  0%,
  100% {
    transform: translateY(0) rotate(-3deg);
  }
  50% {
    transform: translateY(-12px) rotate(3deg);
  }
}

/* ---------- Ряд «заголовок + ссылка» ---------- */
.section-head-row {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
}
.section-head-row :deep(.section-header) {
  margin-bottom: 0;
}

/* ---------- Карточка ценности ---------- */
.value-card {
  padding: var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  height: 100%;
}
.value-card__emoji {
  display: grid;
  place-items: center;
  width: 56px;
  height: 56px;
  font-size: 30px;
  background-color: var(--color-highlight-soft);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-3);
}
.value-card__title {
  font-size: var(--text-lg);
  margin-bottom: var(--space-2);
}
.value-card__text {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}

/* ---------- АДАПТИВ ---------- */
@media (max-width: 860px) {
  .hero__inner {
    grid-template-columns: 1fr;
  }
  /* На мобилке декоративную картинку прячем — она чисто украшение */
  .hero__art {
    display: none;
  }
}
</style>
