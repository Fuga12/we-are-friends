<script setup lang="ts">
import type { Review } from '@/data/types'

defineProps<{ review: Review }>()
</script>

<template>
  <figure class="review-card">
    <div class="review-card__stars" :aria-label="`Оценка ${review.rating} из 5`">
      <span v-for="n in 5" :key="n" aria-hidden="true">{{ n <= review.rating ? '★' : '☆' }}</span>
    </div>
    <blockquote class="review-card__text">«{{ review.text }}»</blockquote>
    <figcaption class="review-card__author">
      <span class="review-card__name">{{ review.author }}</span>
      <span class="review-card__role">{{ review.role }}</span>
    </figcaption>
  </figure>
</template>

<style scoped>
.review-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-6);
  margin: 0;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  height: 100%;
}
.review-card__stars {
  color: var(--brand-amber);
  font-size: var(--text-lg);
  letter-spacing: 2px;
}
.review-card__text {
  margin: 0;
  font-size: var(--text-base);
  line-height: var(--leading-normal);
  color: var(--color-text);
  flex-grow: 1;
}
.review-card__author {
  display: flex;
  flex-direction: column;
  padding-top: var(--space-3);
  border-top: 1px dashed var(--color-border);
}
.review-card__name {
  font-family: var(--font-heading);
  font-weight: var(--weight-bold);
}
.review-card__role {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}
</style>
