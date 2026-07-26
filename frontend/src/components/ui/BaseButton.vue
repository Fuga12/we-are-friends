<script setup lang="ts">
/*
  Универсальная кнопка. Сама решает, чем ей быть:
  - есть `to`   → внутренняя ссылка <router-link> (переход без перезагрузки);
  - есть `href` → внешняя ссылка <a> (например, tel: или vk.com);
  - иначе       → обычная <button>.
  Внешний вид задаёт `variant`, размер — `size`.
*/
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    to?: string
    href?: string
    variant?: 'primary' | 'accent' | 'outline' | 'ghost'
    size?: 'md' | 'lg'
    block?: boolean
    type?: 'button' | 'submit'
  }>(),
  { variant: 'primary', size: 'md', block: false, type: 'button' },
)

// Какой тег рендерим
const tag = computed(() => (props.to ? 'router-link' : props.href ? 'a' : 'button'))

const classes = computed(() => [
  'btn',
  `btn--${props.variant}`,
  `btn--${props.size}`,
  { 'btn--block': props.block },
])
</script>

<template>
  <component
    :is="tag"
    :class="classes"
    :to="to"
    :href="href"
    :type="tag === 'button' ? type : undefined"
  >
    <slot />
  </component>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  line-height: 1;
  text-align: center;
  white-space: nowrap; /* кнопка не переносит текст на две строки */
  border-radius: var(--radius-pill);
  border: 2px solid transparent;
  transition:
    transform var(--transition),
    background-color var(--transition),
    color var(--transition),
    border-color var(--transition),
    box-shadow var(--transition);
}
.btn:active {
  transform: translateY(1px) scale(0.99);
}

/* Размеры */
.btn--md {
  padding: 0.7em 1.4em;
  font-size: var(--text-base);
}
.btn--lg {
  padding: 0.85em 1.9em;
  font-size: var(--text-lg);
}
.btn--block {
  display: flex;
  width: 100%;
}

/* Варианты */
.btn--primary {
  background-color: var(--color-primary);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-md);
}
.btn--primary:hover {
  background-color: var(--color-primary-strong);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-lg);
}

.btn--accent {
  background-color: var(--color-accent);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-md);
}
.btn--accent:hover {
  background-color: var(--color-accent-strong);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-lg);
}

.btn--outline {
  background-color: transparent;
  color: var(--color-primary);
  border-color: var(--color-primary);
}
.btn--outline:hover {
  background-color: var(--color-primary);
  color: var(--color-text-inverse);
}

.btn--ghost {
  background-color: var(--color-surface-muted);
  color: var(--color-text);
}
.btn--ghost:hover {
  background-color: var(--color-border);
  color: var(--color-text);
}
</style>
