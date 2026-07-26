<script setup lang="ts">
/* Яркий блок-призыв «Запишитесь». Используется внизу многих страниц. */
import BaseButton from '@/components/ui/BaseButton.vue'
import { site } from '@/data/site'

withDefaults(
  defineProps<{
    title?: string
    text?: string
  }>(),
  {
    title: 'Готовы познакомиться поближе?',
    text: 'Запишитесь на пробное занятие или задайте вопрос — мы с радостью всё расскажем.',
  },
)
</script>

<template>
  <section class="cta section">
    <div class="container cta__inner">
      <h2 class="cta__title">{{ title }}</h2>
      <p class="cta__text">{{ text }}</p>
      <div class="cta__actions">
        <BaseButton to="/contacts" variant="accent" size="lg">Записаться</BaseButton>
        <BaseButton :href="site.phoneHref" variant="outline" size="lg">
          {{ site.phone }}
        </BaseButton>
      </div>
    </div>
  </section>
</template>

<style scoped>
.cta__inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--space-4);
  padding: clamp(2rem, 1rem + 5vw, 4rem);
  background: linear-gradient(135deg, var(--brand-teal), var(--brand-teal-strong));
  border-radius: var(--radius-xl);
  color: var(--color-text-inverse);
}
.cta__title {
  color: #fff;
  max-width: 20ch;
}
.cta__text {
  max-width: 48ch;
  color: rgba(255, 255, 255, 0.9);
  font-size: var(--text-lg);
}
.cta__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: var(--space-3);
  margin-top: var(--space-2);
}
/* Кнопка-контур на тёмном фоне — делаем белой */
.cta__actions :deep(.btn--outline) {
  color: #fff;
  border-color: rgba(255, 255, 255, 0.7);
}
.cta__actions :deep(.btn--outline:hover) {
  background-color: #fff;
  color: var(--brand-teal-strong);
  border-color: #fff;
}
</style>
