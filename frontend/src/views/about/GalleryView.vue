<script setup lang="ts">
import { computed } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import { useAsync } from '@/composables/useAsync'
import { getGallery } from '@/api'

const { data, loading, error } = useAsync(getGallery)
const gallery = computed(() => data.value ?? [])
</script>

<template>
  <PageHero
    title="Галерея"
    subtitle="Моменты из жизни центра: занятия, мастер-классы, праздники и летние приключения."
    emoji="📸"
    tone="amber"
  />

  <section class="section">
    <div class="container">
      <p v-if="loading" class="text-muted text-center">Загрузка…</p>
      <p v-else-if="error" class="text-muted text-center">
        Не удалось загрузить фотографии. Обновите страницу.
      </p>
      <div v-else class="gallery-grid">
        <figure v-for="item in gallery" :key="item.id" class="tile">
          <img class="tile__img" :src="item.image" alt="Фото центра «Мы Друзья»" loading="lazy" />
        </figure>
      </div>
    </div>
  </section>
</template>

<style scoped>
.gallery-grid {
  display: grid;
  gap: var(--space-4);
  grid-template-columns: repeat(auto-fill, minmax(min(240px, 100%), 1fr));
}
.tile {
  position: relative;
  margin: 0;
  aspect-ratio: 4 / 3;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: transform var(--transition), box-shadow var(--transition);
}
.tile:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}
.tile__img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s var(--ease);
}
.tile:hover .tile__img {
  transform: scale(1.05);
}
</style>
