<script setup lang="ts">
import { watchEffect, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAsync } from '@/composables/useAsync'
import { getNewsPost } from '@/api'
import PhotoFrame from '@/components/ui/PhotoFrame.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const route = useRoute()

// Грузим одну новость по :id из адреса. Второй аргумент заставляет
// перезагрузить данные, если перейти на другую новость без ухода со страницы.
const { data: post, loading } = useAsync(
  () => getNewsPost(route.params.id as string),
  () => route.params.id,
)

// Дата в читаемом виде: «20 июля 2026»
const formattedDate = computed(() =>
  post.value
    ? new Date(post.value.date).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
      })
    : '',
)

watchEffect(() => {
  document.title = post.value ? `${post.value.title} · Мы друзья` : 'Новость не найдена'
})
</script>

<template>
  <div>
    <!-- Обёртка = один стабильный корень страницы. App.vue заворачивает её в
         <Transition name="page">; если менять сам корневой элемент по мере
         загрузки данных (секция «Загрузка» → article), анимация входа
         зависает и страница выглядит пустой. -->
    <section v-if="loading" class="section container text-center">
      <p class="lead">Загрузка…</p>
    </section>
    <section v-else-if="!post" class="section container text-center">
      <p class="lead">Такая новость не найдена.</p>
      <BaseButton to="/news" variant="primary">Ко всем новостям</BaseButton>
    </section>
    <article v-else class="section">
    <div class="container container--narrow">
      <RouterLink to="/news" class="back-link">← Все новости</RouterLink>

      <!-- Мета: категория и дата -->
      <div class="post-meta">
        <span class="post-meta__category">{{ post.category }}</span>
        <time class="post-meta__date" :datetime="post.date">{{ formattedDate }}</time>
      </div>

      <h1 class="post-title">{{ post.title }}</h1>
    </div>

    <!-- Обложка новости -->
    <div class="container container--narrow post-cover">
      <PhotoFrame ratio="16 / 9" :src="post.photo" :alt="post.title" :emoji="post.emoji" label="Фото к новости" />
    </div>

    <!-- Текст новости -->
    <div class="container container--narrow post-body stack">
      <p v-for="(para, i) in post.content" :key="i">{{ para }}</p>
    </div>

    <div class="container container--narrow post-footer">
      <BaseButton to="/news" variant="ghost">← Все новости</BaseButton>
      <BaseButton to="/contacts" variant="accent">Записаться</BaseButton>
    </div>
    </article>
  </div>
</template>

<style scoped>
.back-link {
  display: inline-block;
  margin-bottom: var(--space-5);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--color-primary-strong);
}
.post-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}
.post-meta__category {
  padding: 3px 12px;
  font-size: var(--text-xs);
  font-weight: var(--weight-bold);
  color: var(--color-primary-strong);
  background-color: var(--color-primary-soft);
  border-radius: var(--radius-pill);
}
.post-meta__date {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}
.post-title {
  font-size: var(--text-3xl);
}
.post-cover {
  margin-block: var(--space-6);
}
.post-body {
  font-size: var(--text-lg);
  line-height: var(--leading-normal);
}
.post-body p {
  color: var(--color-text);
}
.post-footer {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
  justify-content: space-between;
  margin-top: var(--space-7);
  padding-top: var(--space-5);
  border-top: 1px solid var(--color-border);
}
</style>
