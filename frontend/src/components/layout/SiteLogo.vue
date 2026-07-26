<script setup lang="ts">
import { site } from '@/data/site'
// Импорт через @/assets — Vite сам подставит правильный путь и захеширует файл.
import logoWordmark from '@/assets/images/logo-wordmark.png'
</script>

<template>
  <RouterLink to="/" class="logo" :aria-label="`${site.name} — детский центр, на главную`">
    <img class="logo__img" :src="logoWordmark" alt="Мы друзья" width="1884" height="292" />
  </RouterLink>
</template>

<style scoped>
.logo {
  display: inline-flex;
  align-items: center;
}
.logo__img {
  /* Высота фиксированная, ширина подстраивается сама (пропорции сохраняются) */
  height: clamp(32px, 2.2vw + 24px, 42px);
  width: auto;
}
</style>
