<script setup lang="ts">
/* Мобильное меню — выезжающая справа панель.
   Пункты с подпунктами раскрываются «гармошкой» по клику. */
import { ref } from 'vue'
import { mainNav } from '@/data/navigation'
import { site } from '@/data/site'
import BaseButton from '@/components/ui/BaseButton.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

// Какие группы сейчас раскрыты (по label). Множество, т.к. можно открыть несколько.
const expanded = ref<Set<string>>(new Set())
function toggle(label: string) {
  const next = new Set(expanded.value)
  next.has(label) ? next.delete(label) : next.add(label)
  expanded.value = next
}
</script>

<template>
  <!-- Teleport переносит меню в конец <body>, минуя шапку.
       Это важно: у шапки есть backdrop-filter, из-за которого fixed-элементы
       внутри неё позиционируются относительно шапки, а не всего экрана. -->
  <Teleport to="body">
    <!-- Затемнение фона. Клик по нему закрывает меню. -->
    <Transition name="fade">
      <div v-if="open" class="backdrop" @click="emit('close')" />
    </Transition>

    <!-- Сама панель -->
    <Transition name="slide">
      <aside v-if="open" class="panel" aria-label="Мобильное меню">
      <nav class="panel__nav">
        <ul>
          <li v-for="link in mainNav" :key="link.to" class="panel__item">
            <!-- Пункт без детей — обычная ссылка -->
            <RouterLink
              v-if="!link.children"
              :to="link.to"
              class="panel__link"
              @click="emit('close')"
            >
              {{ link.label }}
            </RouterLink>

            <!-- Пункт с детьми — кнопка-гармошка -->
            <template v-else>
              <button
                class="panel__link panel__link--toggle"
                :aria-expanded="expanded.has(link.label)"
                @click="toggle(link.label)"
              >
                {{ link.label }}
                <span
                  class="panel__caret"
                  :class="{ 'panel__caret--open': expanded.has(link.label) }"
                  aria-hidden="true"
                  >▾</span
                >
              </button>
              <ul v-show="expanded.has(link.label)" class="panel__sub">
                <li v-for="child in link.children" :key="child.to">
                  <RouterLink :to="child.to" class="panel__sub-link" @click="emit('close')">
                    {{ child.label }}
                  </RouterLink>
                </li>
              </ul>
            </template>
          </li>
        </ul>
      </nav>

      <div class="panel__footer">
        <BaseButton :href="site.phoneHref" variant="primary" block>
          Позвонить
        </BaseButton>
        <a class="panel__phone" :href="site.phoneHref">{{ site.phone }}</a>
      </div>
    </aside>
  </Transition>
  </Teleport>
</template>

<style scoped>
.backdrop {
  position: fixed;
  inset: 0;
  background-color: rgba(43, 38, 34, 0.5);
  z-index: var(--z-overlay);
}
.panel {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: min(88vw, 360px);
  display: flex;
  flex-direction: column;
  padding: calc(var(--header-height) + var(--space-4)) var(--space-5) var(--space-6);
  background-color: var(--color-surface);
  box-shadow: var(--shadow-lg);
  overflow-y: auto;
  z-index: var(--z-menu);
}
.panel__item {
  border-bottom: 1px solid var(--color-border);
}
.panel__link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: var(--space-4) var(--space-1);
  font-family: var(--font-heading);
  font-size: var(--text-lg);
  font-weight: var(--weight-semibold);
  color: var(--color-text);
  text-align: left;
}
.panel__link.router-link-exact-active {
  color: var(--color-primary);
}
.panel__caret {
  transition: transform var(--transition);
  color: var(--color-text-muted);
}
.panel__caret--open {
  transform: rotate(180deg);
}
.panel__sub {
  padding: 0 0 var(--space-3) var(--space-4);
}
.panel__sub-link {
  display: block;
  padding: var(--space-2) var(--space-1);
  font-size: var(--text-base);
  font-weight: var(--weight-medium);
  color: var(--color-text-muted);
}
.panel__sub-link.router-link-exact-active {
  color: var(--color-primary);
}
.panel__footer {
  margin-top: auto;
  padding-top: var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  text-align: center;
}
.panel__phone {
  font-family: var(--font-heading);
  font-weight: var(--weight-bold);
  color: var(--color-text);
}

/* Анимации появления */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition);
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.28s var(--ease);
}
.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}
</style>
