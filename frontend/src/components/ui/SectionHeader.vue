<script setup lang="ts">
/* Единый заголовок секции: маленькая подпись (eyebrow) + заголовок + описание.
   Держит одинаковый ритм на всех страницах. */
withDefaults(
  defineProps<{
    eyebrow?: string
    title: string
    subtitle?: string
    center?: boolean
  }>(),
  { center: false },
)
</script>

<template>
  <header class="section-header" :class="{ 'section-header--center': center }">
    <span v-if="eyebrow" class="eyebrow">{{ eyebrow }}</span>
    <h2>{{ title }}</h2>
    <p v-if="subtitle" class="lead">{{ subtitle }}</p>
  </header>
</template>

<style scoped>
.section-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  margin-bottom: var(--space-6);
  max-width: 62ch;
}
.section-header--center {
  align-items: center;
  text-align: center;
  margin-inline: auto;
}
</style>
