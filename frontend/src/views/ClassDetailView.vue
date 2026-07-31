<script setup lang="ts">
import { watchEffect } from 'vue'
import { useRoute } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getClass } from '@/api'
import { useAccent } from '@/composables/useAccent'
import PhotoFrame from '@/components/ui/PhotoFrame.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import { site } from '@/data/site'

const route = useRoute()

// Грузим занятие по :id из адреса. Второй аргумент перезагружает данные,
// если перейти на другое занятие без ухода со страницы.
const { data: item, loading } = useAsync(
  () => getClass(route.params.id as string),
  () => route.params.id,
)

// Акцентный цвет (HEX) берём у занятия; пока не загрузилось — запасной из useAccent.
const accentStyle = useAccent(() => item.value?.color)

// Заголовок вкладки браузера
watchEffect(() => {
  document.title = item.value ? `${item.value.title} · Мы друзья` : 'Занятие не найдено'
})
</script>

<template>
  <div>
    <!-- Обёртка = один стабильный корень страницы. App.vue заворачивает её в
         <Transition name="page">; если менять сам корневой элемент по мере
         загрузки данных (секция «Загрузка» → article), анимация входа
         зависает на opacity:0 и страница выглядит пустой. -->
    <section v-if="loading" class="section container text-center">
      <p class="lead">Загрузка…</p>
    </section>
    <section v-else-if="!item" class="section container text-center">
      <p class="lead">Такое занятие не найдено.</p>
      <BaseButton to="/classes" variant="primary">Ко всем занятиям</BaseButton>
    </section>
    <article v-else :style="accentStyle">
      <!-- Шапка -->
      <header class="detail-hero">
        <div class="container">
          <RouterLink to="/classes" class="back-link">← Все занятия</RouterLink>
          <div class="detail-hero__top">
            <span class="detail-hero__emoji" aria-hidden="true">{{ item.emoji }}</span>
            <span class="detail-hero__age">{{ item.ageRange }}</span>
          </div>
          <h1 class="detail-hero__title">{{ item.title }}</h1>
          <p class="detail-hero__lead">{{ item.description }}</p>
        </div>
      </header>

      <div class="container detail-body">
        <!-- Левая колонка: фото + описание -->
        <div class="detail-main stack-lg">
          <!-- Главное фото -->
          <PhotoFrame ratio="16 / 9" :src="item.photo" :alt="item.title" :emoji="item.emoji" />

          <!-- Описание -->
          <div class="stack">
            <h2>Об этом занятии</h2>
            <p v-for="(para, i) in item.fullDescription" :key="i">{{ para }}</p>
          </div>

          <!-- Мини-галерея с занятий: список фото приходит с бэкенда -->
          <div v-if="item.photos?.length" class="stack">
            <h2>Фото с занятий</h2>
            <div class="photo-grid">
              <PhotoFrame
                v-for="(src, i) in item.photos"
                :key="src"
                ratio="1 / 1"
                :src="src"
                :alt="`${item.title} — фото ${i + 1}`"
              />
            </div>
          </div>
        </div>

        <!-- Правая колонка: липкая карточка записи -->
        <aside class="detail-aside">
          <div class="info-card">
            <ul v-if="item.highlights" class="info-card__list">
              <li v-for="h in item.highlights" :key="h">{{ h }}</li>
            </ul>
            <dl class="info-card__meta">
              <div>
                <dt>Возраст</dt>
                <dd>{{ item.ageRange }}</dd>
              </div>
              <div>
                <dt>Расписание</dt>
                <dd>{{ item.schedule }}</dd>
              </div>
              <div>
                <dt>Стоимость</dt>
                <dd>{{ item.price }}</dd>
              </div>
            </dl>
            <BaseButton to="/contacts" variant="accent" block>Записаться</BaseButton>
            <BaseButton :href="site.phoneHref" variant="ghost" block>{{ site.phone }}</BaseButton>
          </div>
        </aside>
      </div>
    </article>
  </div>
</template>

<style scoped>
.detail-hero {
  padding-block: var(--space-6) var(--space-7);
  background-color: var(--c-soft);
}
.detail-hero .container {
  position: relative;
}
.back-link {
  display: inline-block;
  margin-bottom: var(--space-4);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--c-strong);
}
.detail-hero__top {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  position: absolute;
  top: var(--space-7);
  right: var(--container-pad);
}
.detail-hero__emoji {
  display: grid;
  place-items: center;
  width: 64px;
  height: 64px;
  font-size: 34px;
  background-color: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.detail-hero__age {
  display: grid;
  place-items: center;
  height: 64px;
  padding-inline: 18px;
  font-weight: var(--weight-semibold);
  color: var(--c-strong);
  background-color: var(--color-surface);
  border-radius: var(--radius-pill);
}
.detail-hero__lead {
  max-width: 60ch;
  margin-top: var(--space-3);
  font-size: var(--text-lg);
  color: var(--color-text-muted);
}

/* Двухколоночная сетка: контент + боковая карточка */
.detail-body {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: var(--space-7);
  padding-block: var(--section-gap);
  align-items: start;
}
.detail-main h2 {
  font-size: var(--text-2xl);
}

.photo-grid {
  display: grid;
  gap: var(--space-3);
  grid-template-columns: repeat(3, 1fr);
}

/* Липкая карточка записи */
.detail-aside {
  position: sticky;
  top: calc(var(--header-height) + var(--space-4));
}
.info-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-top: 5px solid var(--c);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}
.info-card__list {
  display: grid;
  gap: var(--space-2);
}
.info-card__list li {
  position: relative;
  padding-left: 1.6em;
  font-size: var(--text-sm);
}
.info-card__list li::before {
  content: '✓';
  position: absolute;
  left: 0;
  font-weight: var(--weight-bold);
  color: var(--c);
}
.info-card__meta {
  display: grid;
  gap: var(--space-3);
  padding-block: var(--space-4);
  border-block: 1px dashed var(--color-border);
}
.info-card__meta dt {
  font-size: var(--text-xs);
  text-transform: uppercase;
  letter-spacing: 0.03em;
  color: var(--color-text-muted);
}
.info-card__meta dd {
  font-weight: var(--weight-bold);
}

/* Адаптив: на планшете и мобилке — одна колонка, карточка не липкая */
@media (max-width: 860px) {
  .detail-hero__top {
    position: static;
    margin-bottom: var(--space-3);
  }
  .detail-body {
    grid-template-columns: 1fr;
  }
  .detail-aside {
    position: static;
  }
}
</style>
