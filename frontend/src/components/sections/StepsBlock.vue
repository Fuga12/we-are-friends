<script setup lang="ts">
/* Пронумерованные шаги «как это проходит». Переиспользуется на спец-страницах. */
defineProps<{
  title?: string
  steps: { title: string; text: string }[]
}>()
</script>

<template>
  <section class="section">
    <div class="container">
      <h2 v-if="title" class="steps-title">{{ title }}</h2>
      <ol class="steps">
        <li v-for="(step, i) in steps" :key="i" class="step">
          <span class="step__num">{{ i + 1 }}</span>
          <h3 class="step__title">{{ step.title }}</h3>
          <p class="step__text">{{ step.text }}</p>
        </li>
      </ol>
    </div>
  </section>
</template>

<style scoped>
.steps-title {
  margin-bottom: var(--space-6);
  text-align: center;
}
.steps {
  display: grid;
  gap: var(--space-5);
  grid-template-columns: repeat(auto-fit, minmax(min(220px, 100%), 1fr));
  counter-reset: step;
}
.step {
  padding: var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.step__num {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  margin-bottom: var(--space-3);
  font-family: var(--font-heading);
  font-weight: var(--weight-bold);
  font-size: var(--text-lg);
  color: var(--color-text-inverse);
  background-color: var(--color-accent);
  border-radius: 50%;
}
.step__title {
  font-size: var(--text-lg);
  margin-bottom: var(--space-2);
}
.step__text {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
</style>
