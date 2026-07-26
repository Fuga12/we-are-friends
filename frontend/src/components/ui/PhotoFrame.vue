<script setup lang="ts">
/*
  Область под фотографию.
  - Если передать `src` — покажет настоящее фото (обрезав по размеру, object-fit: cover).
  - Если `src` нет — покажет аккуратную заглушку с подсказкой «Фото».

  Так весь сайт уже готов к реальным фото: когда они появятся, достаточно
  передать в этот компонент ссылку на картинку, и плейсхолдер сам заменится.

  `ratio` — пропорции области (например '4 / 3', '16 / 9', '1 / 1').
  Цвет подложки берётся из переменной --c-soft, если родитель её задал
  (см. useAccent), иначе — нейтральный.
*/
withDefaults(
  defineProps<{
    src?: string
    alt?: string
    emoji?: string
    ratio?: string
    label?: string
  }>(),
  { ratio: '4 / 3', alt: '', label: 'Фото' },
)
</script>

<template>
  <div class="photo" :style="{ aspectRatio: ratio }">
    <img v-if="src" :src="src" :alt="alt" class="photo__img" loading="lazy" />

    <div v-else class="photo__placeholder">
      <span v-if="emoji" class="photo__emoji" aria-hidden="true">{{ emoji }}</span>
      <span class="photo__hint">
        <span aria-hidden="true">📷</span>
        {{ label }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.photo {
  position: relative;
  width: 100%;
  overflow: hidden;
  border-radius: var(--radius-md);
  background-color: var(--c-soft, var(--color-surface-muted));
}
.photo__img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.photo__placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  /* Пунктирная рамка — визуальный сигнал «сюда встанет фото» */
  border: 2px dashed color-mix(in srgb, var(--c, var(--color-text-muted)) 35%, transparent);
  border-radius: inherit;
}
.photo__emoji {
  font-size: clamp(32px, 6vw, 48px);
  opacity: 0.85;
}
.photo__hint {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 12px;
  font-family: var(--font-heading);
  font-size: var(--text-xs);
  font-weight: var(--weight-semibold);
  color: var(--color-text-muted);
  background-color: rgba(255, 255, 255, 0.7);
  border-radius: var(--radius-pill);
}
</style>
