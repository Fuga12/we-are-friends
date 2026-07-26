<script setup lang="ts">
/* Каркас раздела «О центре»: под-навигация + <router-view> для подстраниц.
   Дочерние маршруты (about/gallery, about/staff и т.д.) рендерятся внутрь. */
import { aboutNav } from '@/data/navigation'
</script>

<template>
  <div>
    <!-- Липкая под-навигация раздела, прокручивается по горизонтали на мобилке -->
    <nav class="subnav" aria-label="Навигация по разделу «О центре»">
      <div class="container subnav__inner">
        <RouterLink
          v-for="link in aboutNav"
          :key="link.to"
          :to="link.to"
          class="subnav__link"
          exact-active-class="subnav__link--active"
        >
          {{ link.label }}
        </RouterLink>
      </div>
    </nav>

    <RouterView />
  </div>
</template>

<style scoped>
.subnav {
  position: sticky;
  top: var(--header-height);
  z-index: 50;
  background-color: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
}
.subnav__inner {
  display: flex;
  gap: var(--space-1);
  /* Горизонтальная прокрутка, если ссылки не помещаются (мобилка) */
  overflow-x: auto;
  scrollbar-width: none; /* прячем полосу прокрутки */
}
.subnav__inner::-webkit-scrollbar {
  display: none;
}
.subnav__link {
  flex-shrink: 0;
  padding: var(--space-4) var(--space-3);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--color-text-muted);
  border-bottom: 3px solid transparent;
  white-space: nowrap;
  transition: color var(--transition), border-color var(--transition);
}
.subnav__link:hover {
  color: var(--color-text);
}
.subnav__link--active {
  color: var(--color-primary-strong);
  border-bottom-color: var(--color-primary);
}
</style>
