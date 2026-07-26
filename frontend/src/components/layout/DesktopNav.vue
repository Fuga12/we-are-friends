<script setup lang="ts">
/* Горизонтальное меню для десктопа. У пунктов с children появляется
   выпадающий список по наведению И по фокусу с клавиатуры (доступность). */
import { mainNav } from '@/data/navigation'
</script>

<template>
  <nav class="desktop-nav" aria-label="Основное меню">
    <ul class="desktop-nav__list">
      <li
        v-for="link in mainNav"
        :key="link.to"
        class="desktop-nav__item"
        :class="{ 'has-children': link.children }"
      >
        <RouterLink :to="link.to" class="desktop-nav__link">
          {{ link.label }}
          <span v-if="link.children" class="desktop-nav__caret" aria-hidden="true">▾</span>
        </RouterLink>

        <!-- Выпадающее подменю -->
        <ul v-if="link.children" class="desktop-nav__dropdown">
          <li v-for="child in link.children" :key="child.to">
            <RouterLink :to="child.to" class="desktop-nav__dropdown-link">
              {{ child.label }}
            </RouterLink>
          </li>
        </ul>
      </li>
    </ul>
  </nav>
</template>

<style scoped>
.desktop-nav__list {
  display: flex;
  align-items: center;
  gap: var(--space-1);
}
.desktop-nav__item {
  position: relative;
}
.desktop-nav__link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: var(--space-2) var(--space-3);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--color-text);
  white-space: nowrap; /* «О центре» и др. не переносятся на две строки */
  border-radius: var(--radius-pill);
  transition: background-color var(--transition), color var(--transition);
}
.desktop-nav__link:hover {
  color: var(--color-primary-strong);
  background-color: var(--color-primary-soft);
}
/* Активная ссылка текущей страницы (класс от <router-link>) */
.desktop-nav__link.router-link-active {
  color: var(--color-primary-strong);
}
.desktop-nav__caret {
  font-size: 10px;
  transition: transform var(--transition);
}

/* --- Выпадающее меню --- */
.desktop-nav__dropdown {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  min-width: 210px;
  padding: var(--space-2);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  /* Скрыто по умолчанию, показываем плавно */
  opacity: 0;
  visibility: hidden;
  transform: translateY(6px);
  transition:
    opacity var(--transition),
    transform var(--transition),
    visibility var(--transition);
  z-index: var(--z-header);
}
/* Показать при наведении на пункт или при фокусе внутри него */
.has-children:hover .desktop-nav__dropdown,
.has-children:focus-within .desktop-nav__dropdown {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}
.has-children:hover .desktop-nav__caret {
  transform: rotate(180deg);
}
.desktop-nav__dropdown-link {
  display: block;
  padding: var(--space-2) var(--space-3);
  font-size: var(--text-sm);
  font-weight: var(--weight-medium);
  color: var(--color-text);
  border-radius: var(--radius-sm);
}
.desktop-nav__dropdown-link:hover,
.desktop-nav__dropdown-link.router-link-exact-active {
  background-color: var(--color-primary-soft);
  color: var(--color-primary-strong);
}
</style>
