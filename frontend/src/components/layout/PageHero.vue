<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    emoji?: string
    tone?: 'teal' | 'amber' | 'orange'
  }>(),
  { tone: 'teal' },
)
</script>

<template>
  <section class="page-hero" :class="`page-hero--${tone}`">
    <div class="container page-hero__inner">
      <!-- <span v-if="emoji" class="page-hero__emoji" aria-hidden="true">{{ emoji }}</span> -->
      <h1 class="page-hero__title">{{ title }}</h1>
      <p v-if="subtitle" class="page-hero__subtitle">{{ subtitle }}</p>
    </div>
  </section>
</template>

<style scoped>
.page-hero {
  padding-block: clamp(2.5rem, 1.5rem + 4vw, 4.5rem);
  text-align: center;
}
.page-hero--teal {
  background-color: var(--color-primary-soft);
}
.page-hero--amber {
  background-color: var(--color-highlight-soft);
}
.page-hero--orange {
  background-color: var(--color-accent-soft);
}
.page-hero__inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
}
.page-hero__emoji {
  font-size: clamp(2.5rem, 2rem + 3vw, 4rem);
  line-height: 1;
}
.page-hero__subtitle {
  max-width: 55ch;
  font-size: var(--text-lg);
  color: var(--color-text-muted);
}
</style>
