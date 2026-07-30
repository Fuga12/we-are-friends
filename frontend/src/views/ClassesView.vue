<script setup lang="ts">
import { computed } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import ClassCard from '@/components/cards/ClassCard.vue'
import CtaSection from '@/components/sections/CtaSection.vue'
import { useAsync } from '@/composables/useAsync'
import { getClasses } from '@/api'

const { data, loading, error } = useAsync(getClasses)
const classes = computed(() => data.value ?? [])
</script>

<template>
  <PageHero
    title="Занятия и направления"
    subtitle="От первых шагов в развитии до подготовки к школе — подберём то, что подойдёт именно вашему ребёнку."
    emoji="🎈"
  />

  <section class="section">
    <div class="container">
      <p v-if="loading" class="text-muted text-center">Загрузка…</p>
      <p v-else-if="error" class="text-muted text-center">
        Не удалось загрузить занятия. Обновите страницу.
      </p>
      <div v-else class="grid" style="--min: 300px">
        <ClassCard v-for="item in classes" :key="item.id" :item="item" />
      </div>

      <p class="text-muted text-center classes-note">
        Не нашли подходящее направление? Позвоните нам — подскажем и поможем с выбором.
      </p>
    </div>
  </section>

  <CtaSection
    title="Приходите на пробное занятие"
    text="Первое знакомство — бесплатно. Ребёнок попробует, а вы посмотрите, как всё устроено."
  />
</template>

<style scoped>
.classes-note {
  margin-top: var(--space-7);
}
</style>
