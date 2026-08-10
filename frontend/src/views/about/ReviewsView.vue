<script setup lang="ts">
import { computed } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import ReviewCard from '@/components/cards/ReviewCard.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import { useAsync } from '@/composables/useAsync'
import { getReviews } from '@/api'
import { site } from '@/data/site'
const {data,loading,error}=useAsync(getReviews)
const reviews=computed(()=>data.value??[])
</script>

<template>
  <PageHero
    title="Отзывы родителей"
    subtitle="Спасибо, что делитесь впечатлениями — это лучшая награда для нашей команды."
    emoji="💬"
    tone="orange"
  />

  <section class="section">
    <div class="container">
      <p v-if="loading" class="text-muted text-center">Загрузка…</p>
      <p v-else-if="error" class="text-muted text-center">Не удалось загрузить отзывы. Обновите страницу.</p>
      <div class="grid" style="--min: 320px">
        <ReviewCard v-for="review in reviews" :key="review.id" :review="review" />
      </div>

      <div class="reviews-cta">
        <p class="lead">Были у нас? Оставьте отзыв — это помогает другим родителям выбрать.</p>
        <BaseButton :href="site.vk" variant="primary" size="lg">
          Оставить отзыв во ВКонтакте
        </BaseButton>
      </div>
    </div>
  </section>
</template>

<style scoped>
.reviews-cta {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  text-align: center;
  margin-top: var(--space-8);
  padding: var(--space-7) var(--space-5);
  background-color: var(--color-highlight-soft);
  border-radius: var(--radius-xl);
}
</style>
