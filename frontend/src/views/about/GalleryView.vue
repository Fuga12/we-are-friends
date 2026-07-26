<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import { gallery } from '@/data/gallery'
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
      <div class="gallery-grid">
        <figure v-for="item in gallery" :key="item.id" class="tile">
          <img class="tile__img" :src="item.image" :alt="item.title ?? 'Фото центра'" loading="lazy" />
          <figcaption v-if="item.title" class="tile__caption">{{ item.title }}</figcaption>
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
.tile__caption {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  padding: var(--space-5) var(--space-4) var(--space-3);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: #fff;
  /* мягкая тёмная подложка снизу, чтобы подпись читалась на любом фото */
  background: linear-gradient(to top, rgba(43, 38, 34, 0.75), transparent);
}
</style>
