<script setup lang="ts">
import { toRef } from 'vue'
import { useAccent } from '@/composables/useAccent'
import type { StaffMember } from '@/data/types'

const props = defineProps<{ member: StaffMember }>()
const accentStyle = useAccent(toRef(() => props.member.color))
</script>

<template>
  <article class="staff-card" :style="accentStyle">
    <!-- Заглушка-аватар с инициалами. Позже сюда можно вставить <img>. -->
    <div class="staff-card__avatar" aria-hidden="true">{{ member.initials }}</div>
    <h3 class="staff-card__name">{{ member.name }}</h3>
    <p class="staff-card__role">{{ member.role }}</p>
    <p class="staff-card__bio">{{ member.bio }}</p>
  </article>
</template>

<style scoped>
.staff-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--space-2);
  padding: var(--space-6) var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  height: 100%;
  transition:
    transform var(--transition),
    box-shadow var(--transition);
}
.staff-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}
.staff-card__avatar {
  display: grid;
  place-items: center;
  width: 92px;
  height: 92px;
  margin-bottom: var(--space-2);
  font-family: var(--font-heading);
  font-size: var(--text-2xl);
  font-weight: var(--weight-bold);
  color: var(--color-text-inverse);
  background: linear-gradient(135deg, var(--c), var(--c-strong));
  border-radius: 50%;
  box-shadow: var(--shadow-md);
}
.staff-card__name {
  font-size: var(--text-lg);
}
.staff-card__role {
  font-weight: var(--weight-bold);
  color: var(--c-strong);
  font-size: var(--text-sm);
}
.staff-card__bio {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
</style>
