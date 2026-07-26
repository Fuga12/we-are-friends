<script setup lang="ts">
import { computed } from 'vue'
import type { NewsPost } from '@/data/types'

const props = defineProps<{ post: NewsPost }>()

// Форматируем ISO-дату в читаемый русский вид: «20 июля 2026»
const formattedDate = computed(() =>
  new Date(props.post.date).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }),
)
</script>

<template>
  <RouterLink class="news-card" :to="`/news/${post.id}`">
    <img
      v-if="post.photo"
      class="news-card__photo"
      :src="post.photo"
      :alt="post.title"
      loading="lazy"
    />
    <div v-else class="news-card__emoji" aria-hidden="true">{{ post.emoji }}</div>
    <div class="news-card__body">
      <div class="news-card__top">
        <span class="news-card__category">{{ post.category }}</span>
        <time class="news-card__date" :datetime="post.date">{{ formattedDate }}</time>
      </div>
      <h3 class="news-card__title">{{ post.title }}</h3>
      <p class="news-card__excerpt">{{ post.excerpt }}</p>
      <span class="news-card__more">Читать →</span>
    </div>
  </RouterLink>
</template>

<style scoped>
.news-card {
  display: flex;
  gap: var(--space-4);
  padding: var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  color: var(--color-text); /* карточка — ссылка, отменяем цвет ссылки */
  transition:
    transform var(--transition),
    box-shadow var(--transition);
}
.news-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  color: var(--color-text);
}
.news-card:hover .news-card__title {
  color: var(--color-primary-strong);
}
.news-card__emoji {
  display: grid;
  place-items: center;
  flex-shrink: 0;
  width: 56px;
  height: 56px;
  font-size: 28px;
  background-color: var(--color-highlight-soft);
  border-radius: var(--radius-md);
}
.news-card__photo {
  flex-shrink: 0;
  width: 96px;
  height: 96px;
  object-fit: cover;
  border-radius: var(--radius-md);
}
.news-card__body {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}
.news-card__top {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: var(--space-3);
}
.news-card__category {
  padding: 2px 10px;
  font-size: var(--text-xs);
  font-weight: var(--weight-bold);
  color: var(--color-primary-strong);
  background-color: var(--color-primary-soft);
  border-radius: var(--radius-pill);
}
.news-card__date {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}
.news-card__title {
  font-size: var(--text-lg);
}
.news-card__excerpt {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
.news-card__more {
  align-self: flex-start;
  margin-top: var(--space-1);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--color-primary-strong);
}

/* На узких экранах эмодзи-плитка уезжает наверх, текст под ней */
@media (max-width: 480px) {
  .news-card {
    flex-direction: column;
  }
}
</style>
