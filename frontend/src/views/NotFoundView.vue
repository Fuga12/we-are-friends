<script setup lang="ts">
import BaseButton from '@/components/ui/BaseButton.vue'
</script>

<template>
  <section class="notfound section">
    <div class="container">
      <span class="notfound__emoji" aria-hidden="true">🧩</span>
      <h1 class="notfound__title">Упс! Страница потерялась</h1>
      <p class="notfound__text">
        Кажется, такой страницы у нас нет. Возможно, ссылка устарела или в адресе опечатка.
      </p>
      <BaseButton to="/" variant="primary" size="lg">Вернуться на главную</BaseButton>
    </div>
  </section>
</template>

<style scoped>
.notfound {
  text-align: center;
}
.notfound .container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
}
.notfound__emoji {
  font-size: clamp(64px, 12vw, 110px);
}
.notfound__text {
  max-width: 44ch;
  color: var(--color-text-muted);
  font-size: var(--text-lg);
}
</style>
