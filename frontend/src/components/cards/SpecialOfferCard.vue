<script setup lang="ts">
import { toRef } from 'vue'
import { useAccent } from '@/composables/useAccent'
import BaseButton from '@/components/ui/BaseButton.vue'
import type { SpecialOffer } from '@/data/types'

const props = defineProps<{ offer: SpecialOffer }>()
const accentStyle = useAccent(toRef(() => props.offer.color))
</script>

<template>
  <article class="offer-card" :style="accentStyle">
    <div class="offer-card__head">
      <span class="offer-card__emoji" aria-hidden="true">{{ offer.emoji }}</span>
      <div>
        <span class="offer-card__tagline">{{ offer.tagline }}</span>
        <h3 class="offer-card__title">{{ offer.title }}</h3>
      </div>
    </div>
    <p class="offer-card__desc">{{ offer.description }}</p>
    <ul class="offer-card__features">
      <li v-for="feature in offer.features" :key="feature">{{ feature }}</li>
    </ul>
    <div class="offer-card__footer">
      <span class="offer-card__price">{{ offer.price }}</span>
      <BaseButton :to="offer.to" variant="outline">Подробнее</BaseButton>
    </div>
  </article>
</template>

<style scoped>
.offer-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-6);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-md);
  height: 100%;
  transition:
    transform var(--transition),
    box-shadow var(--transition);
}
.offer-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-lg);
}
.offer-card__head {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}
.offer-card__emoji {
  display: grid;
  place-items: center;
  flex-shrink: 0;
  width: 64px;
  height: 64px;
  font-size: 34px;
  background-color: var(--c-soft);
  border-radius: var(--radius-lg);
}
.offer-card__tagline {
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--c-strong);
}
.offer-card__title {
  font-size: var(--text-xl);
}
.offer-card__desc {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
.offer-card__features {
  display: grid;
  gap: var(--space-2);
  flex-grow: 1;
}
.offer-card__features li {
  position: relative;
  padding-left: 1.6em;
  font-size: var(--text-sm);
}
.offer-card__features li::before {
  content: '✓';
  position: absolute;
  left: 0;
  font-weight: var(--weight-bold);
  color: var(--c);
}
.offer-card__footer {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px dashed var(--color-border);
}
.offer-card__price {
  font-family: var(--font-heading);
  font-size: var(--text-xl);
  font-weight: var(--weight-bold);
  color: var(--c-strong);
}
</style>
