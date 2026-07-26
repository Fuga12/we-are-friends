<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import StaffCard from '@/components/cards/StaffCard.vue'
import CtaSection from '@/components/sections/CtaSection.vue'
import { staff } from '@/data/staff'
</script>

<template>
  <PageHero
    title="Наша команда"
    subtitle="Педагоги и специалисты, которым можно доверить самое дорогое."
    emoji="👩‍🏫"
  />

  <section class="section">
    <div class="container">
      <div class="grid" style="--min: 240px">
        <StaffCard v-for="member in staff" :key="member.id" :member="member" />
      </div>
    </div>
  </section>

  <CtaSection
    title="Хотите присоединиться к команде?"
    text="Мы всегда рады неравнодушным педагогам. Напишите нам — расскажем об открытых вакансиях."
  />
</template>
