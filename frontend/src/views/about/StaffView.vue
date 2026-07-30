<script setup lang="ts">
import { computed } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import StaffCard from '@/components/cards/StaffCard.vue'
import CtaSection from '@/components/sections/CtaSection.vue'
import { useAsync } from '@/composables/useAsync'
import { getStaff } from '@/api'

const { data, loading, error } = useAsync(getStaff)
const staff = computed(() => data.value ?? [])
</script>

<template>
  <PageHero
    title="Наша команда"
    subtitle="Педагоги и специалисты, которым можно доверить самое дорогое."
    emoji="👩‍🏫"
  />

  <section class="section">
    <div class="container">
      <p v-if="loading" class="text-muted text-center">Загрузка…</p>
      <p v-else-if="error" class="text-muted text-center">
        Не удалось загрузить сотрудников. Обновите страницу.
      </p>
      <div v-else class="grid" style="--min: 240px">
        <StaffCard v-for="member in staff" :key="member.id" :member="member" />
      </div>
    </div>
  </section>

  <CtaSection
    title="Хотите присоединиться к команде?"
    text="Мы всегда рады неравнодушным педагогам. Напишите нам — расскажем об открытых вакансиях."
  />
</template>
