<script setup lang="ts">
import { computed, ref } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import NewsCard from '@/components/cards/NewsCard.vue'
import { useAsync } from '@/composables/useAsync'
import { getNews } from '@/api'

const { data, loading, error } = useAsync(getNews)
const news = computed(() => data.value ?? [])

// Собираем список категорий из данных (уникальные значения)
const categories = computed(() => ['Все', ...new Set(news.value.map((n) => n.category))])
const active = ref('Все')

const visible = computed(() =>
  active.value === 'Все' ? news.value : news.value.filter((n) => n.category === active.value),
)
</script>

<template>
  <PageHero
    title="Новости центра"
    subtitle="Анонсы, события и достижения наших воспитанников."
    emoji="📰"
    tone="amber"
  />

  <section class="section">
    <div class="container container--narrow">
      <p v-if="loading" class="text-muted text-center">Загрузка…</p>
      <p v-else-if="error" class="text-muted text-center">
        Не удалось загрузить новости. Обновите страницу.
      </p>
      <template v-else>
        <!-- Фильтр по категориям -->
        <div class="chips">
          <button
            v-for="cat in categories"
            :key="cat"
            class="chip"
            :class="{ 'chip--active': active === cat }"
            @click="active = cat"
          >
            {{ cat }}
          </button>
        </div>

        <div class="stack-lg">
          <NewsCard v-for="post in visible" :key="post.id" :post="post" />
        </div>
      </template>
    </div>
  </section>
</template>

<style scoped>
.chips {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
  margin-bottom: var(--space-6);
}
.chip {
  padding: var(--space-2) var(--space-4);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--color-text-muted);
  background-color: var(--color-surface-muted);
  border-radius: var(--radius-pill);
  transition: all var(--transition);
}
.chip--active {
  color: var(--color-text-inverse);
  background-color: var(--color-primary);
}
</style>
