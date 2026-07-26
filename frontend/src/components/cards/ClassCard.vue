<script setup lang="ts">
import { toRef } from 'vue'
import { useAccent } from '@/composables/useAccent'
import PhotoFrame from '@/components/ui/PhotoFrame.vue'
import type { ClassItem } from '@/data/types'

const props = defineProps<{ item: ClassItem }>()
const accentStyle = useAccent(toRef(() => props.item.color))
</script>

<template>
  <RouterLink class="class-card" :to="`/classes/${item.id}`" :style="accentStyle">
    <!-- Область под фото занятия. Пока заглушка с эмодзи; когда появится
         реальное фото, передай в PhotoFrame :src="item.photo". -->
    <div class="class-card__media">
      <PhotoFrame ratio="16 / 10" :src="item.photo" :alt="item.title" :emoji="item.emoji" />
      <span class="class-card__age">{{ item.ageRange }}</span>
    </div>

    <div class="class-card__body">
      <h3 class="class-card__title">{{ item.title }}</h3>
      <p class="class-card__desc">{{ item.description }}</p>
      <dl class="class-card__meta">
        <div>
          <dt>Расписание</dt>
          <dd>{{ item.schedule }}</dd>
        </div>
        <div>
          <dt>Стоимость</dt>
          <dd>{{ item.price }}</dd>
        </div>
      </dl>
      <span class="class-card__more">Подробнее →</span>
    </div>
  </RouterLink>
</template>

<style scoped>
.class-card {
  display: flex;
  flex-direction: column;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-top: 5px solid var(--c);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  overflow: hidden; /* чтобы фото не выходило за скруглённые углы */
  color: var(--color-text); /* карточка — ссылка, отменяем цвет ссылки */
  transition:
    transform var(--transition),
    box-shadow var(--transition);
  height: 100%;
}
.class-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-lg);
  color: var(--color-text);
}
.class-card:hover .class-card__title {
  color: var(--c-strong);
}
.class-card:hover .class-card__more {
  gap: var(--space-2);
}

/* Фото-шапка карточки */
.class-card__media {
  position: relative;
}
.class-card__media :deep(.photo) {
  border-radius: 0; /* верх карточки — фото во всю ширину, без скруглений */
}
.class-card__age {
  position: absolute;
  top: var(--space-3);
  left: var(--space-3);
  padding: 3px 12px;
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--c-strong);
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: var(--radius-pill);
  box-shadow: var(--shadow-sm);
}

/* Текстовая часть */
.class-card__body {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-5);
  flex-grow: 1;
}
.class-card__title {
  font-size: var(--text-xl);
}
.class-card__desc {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
  flex-grow: 1;
}
.class-card__meta {
  display: grid;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px dashed var(--color-border);
}
.class-card__meta dt {
  font-size: var(--text-xs);
  text-transform: uppercase;
  letter-spacing: 0.03em;
  color: var(--color-text-muted);
}
.class-card__meta dd {
  font-weight: var(--weight-bold);
  color: var(--color-text);
}
.class-card__more {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  margin-top: var(--space-1);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
  color: var(--c-strong);
  transition: gap var(--transition);
}
</style>
