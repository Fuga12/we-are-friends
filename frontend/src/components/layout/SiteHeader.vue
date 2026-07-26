<script setup lang="ts">
import { ref, watch, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import SiteLogo from './SiteLogo.vue'
import DesktopNav from './DesktopNav.vue'
import MobileMenu from './MobileMenu.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import { site } from '@/data/site'

const route = useRoute()
const menuOpen = ref(false)

// Блокируем прокрутку страницы, пока открыто мобильное меню
watch(menuOpen, (open) => {
  document.body.style.overflow = open ? 'hidden' : ''
})

// Закрываем меню при переходе на другую страницу
watch(
  () => route.fullPath,
  () => {
    menuOpen.value = false
  },
)

// Закрытие по клавише Esc
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') menuOpen.value = false
}
window.addEventListener('keydown', onKeydown)
onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>

<template>
  <header class="header">
    <div class="header__inner container">
      <SiteLogo />

      <!-- Десктопное меню (скрыто на мобилке через CSS) -->
      <DesktopNav class="header__desktop-nav" />

      <div class="header__actions">
        <BaseButton :href="site.phoneHref" variant="accent" class="header__cta">
          {{ site.phone }}
        </BaseButton>

        <!-- Кнопка-бургер (видна только на мобилке) -->
        <button
          class="burger"
          :class="{ 'burger--open': menuOpen }"
          :aria-expanded="menuOpen"
          aria-label="Открыть меню"
          @click="menuOpen = !menuOpen"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>
    </div>

    <MobileMenu :open="menuOpen" @close="menuOpen = false" />
  </header>
</template>

<style scoped>
.header {
  position: sticky;
  top: 0;
  z-index: var(--z-header);
  background-color: rgba(255, 253, 249, 0.85);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--color-border);
}
.header__inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  min-height: var(--header-height);
}
.header__actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

/* Бургер: три полоски, превращаются в крестик */
.burger {
  display: none; /* показываем только на мобилке (см. медиа-запрос) */
  flex-direction: column;
  justify-content: center;
  gap: 5px;
  width: 44px;
  height: 44px;
  padding: 10px;
  border-radius: var(--radius-md);
}
.burger span {
  display: block;
  height: 3px;
  width: 100%;
  background-color: var(--color-text);
  border-radius: 2px;
  transition: transform var(--transition), opacity var(--transition);
}
.burger--open span:nth-child(1) {
  transform: translateY(8px) rotate(45deg);
}
.burger--open span:nth-child(2) {
  opacity: 0;
}
.burger--open span:nth-child(3) {
  transform: translateY(-8px) rotate(-45deg);
}

/* ============ АДАПТИВ ============
   Точка перелома 1100px: ниже неё горизонтальное меню и телефон-кнопка
   уже не помещаются в одну строку и «съезжают», поэтому прячем их и
   показываем бургер. */
@media (max-width: 1100px) {
  .header__desktop-nav {
    display: none;
  }
  .header__cta {
    display: none;
  }
  .burger {
    display: flex;
  }
}
</style>
