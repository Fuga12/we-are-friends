<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import ClassCard from '@/components/cards/ClassCard.vue'
import CtaSection from '@/components/sections/CtaSection.vue'
import { classes } from '@/data/classes'
</script>

<template>
  <PageHero
    title="Занятия и направления"
    subtitle="От первых шагов в развитии до подготовки к школе — подберём то, что подойдёт именно вашему ребёнку."
    emoji="🎈"
  />

  <section class="section">
    <div class="container">
      <div class="grid" style="--min: 300px">
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
