<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import StepsBlock from '@/components/sections/StepsBlock.vue'
import CtaSection from '@/components/sections/CtaSection.vue'

const packages = [
  {
    emoji: '🧁',
    title: 'Кулинарный праздник',
    price: '8 000 ₽',
    text: 'Аренда, чаепитие и кулинарный мастер-класс (1–1,5 ч) на выбор: вафли, пончики, капкейки, калитки, пицца или бенто-торт. До 10 детей.',
  },
  {
    emoji: '🪚',
    title: 'Столярный праздник',
    price: '8 000 ₽',
    text: 'Аренда, чаепитие и столярный мастер-класс на выбор: нож-кунай, змейка или подсвечник «Крипер». Каждый унесёт изделие домой. До 10 детей.',
  },
  {
    emoji: '✨',
    title: 'Дополнительно',
    price: '+ к пакету',
    text: 'Анимация (Человек-паук, Леди Баг и Супер Кот или интерактив) +4 000 ₽. Серебряная дискотека +1 500 ₽.',
  },
]

const steps = [
  { title: 'Выбираете дату', text: 'Звоните нам, подбираем свободное время и обсуждаем формат.' },
  { title: 'Выбираете мастер-класс', text: 'Кулинарный или столярный — и конкретное изделие или блюдо.' },
  { title: 'Мы всё готовим', text: 'Помещение, чаепитие и материалы для мастер-класса. Торт и сладости — от вас.' },
  { title: 'Праздник!', text: 'Дети мастерят, играют и веселятся, а вы просто наслаждаетесь моментом.' },
]
</script>

<template>
  <PageHero
    title="День рождения ребёнка"
    subtitle="Организуем праздник «под ключ», чтобы у ребёнка был самый счастливый день, а у вас — ни капли хлопот."
    emoji="🎉"
    tone="orange"
  />

  <section class="section">
    <div class="container">
      <SectionHeader eyebrow="Пакеты" title="Выберите формат праздника" center />
      <div class="grid" style="--min: 260px">
        <article v-for="p in packages" :key="p.title" class="pkg">
          <span class="pkg__emoji" aria-hidden="true">{{ p.emoji }}</span>
          <h3 class="pkg__title">{{ p.title }}</h3>
          <span class="pkg__price">{{ p.price }}</span>
          <p class="pkg__text">{{ p.text }}</p>
        </article>
      </div>
    </div>
  </section>

  <div class="section--muted">
    <StepsBlock title="Как всё проходит" :steps="steps" />
  </div>

  <CtaSection
    title="Забронируйте дату праздника"
    text="Популярные выходные разбирают заранее — позвоните, чтобы выбрать удобное время."
  />
</template>

<style scoped>
.pkg {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--space-2);
  padding: var(--space-6);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.pkg__emoji {
  font-size: 44px;
}
.pkg__title {
  font-size: var(--text-lg);
}
.pkg__price {
  font-family: var(--font-heading);
  font-weight: var(--weight-bold);
  color: var(--color-accent-strong);
  font-size: var(--text-lg);
}
.pkg__text {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
</style>
