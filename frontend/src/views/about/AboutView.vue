<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import PhotoFrame from '@/components/ui/PhotoFrame.vue'
import CtaSection from '@/components/sections/CtaSection.vue'
import { site } from '@/data/site'

const facts = [
  { emoji: '🎨', text: '15 направлений: от столярки и шитья до робототехники и кулинарии.' },
  { emoji: '🏠', text: 'Две площадки в Петрозаводске: пр. А. Невского, 60 и филиал на Древлянке.' },
  { emoji: '🏆', text: 'Участвуем и побеждаем в конкурсах, фестивалях и модных показах.' },
  { emoji: '☀️', text: 'Летние программы: «Школа магии», «Деффчонки», «Майнкрафт».' },
]

const principles = [
  { emoji: '❤️', title: 'Мы — семья', text: 'Тёплое и уважительное отношение к каждому ребёнку, а не конвейер.' },
  { emoji: '🛠️', title: 'Всё по-настоящему', text: 'Настоящие инструменты и материалы. Дети создают реальные вещи.' },
  { emoji: '🎯', title: 'Больше, чем ремесло', text: 'Учим общаться, проявлять себя, ставить и добиваться целей.' },
]
</script>

<template>
  <PageHero
    title="О центре «Мы Друзья»"
    subtitle="Площадка для творческого развития детей и взрослых в Петрозаводске."
    emoji="🏠"
  />

  <!-- Вступление -->
  <section class="section">
    <div class="container container--narrow stack-lg">
      <div class="stack">
        <p class="lead">
          {{ site.name }} — это не просто место, где дети учатся и мастерят. Это семья! Небольшой
          центр с тёплым и уважительным отношением к каждому ребёнку.
        </p>
        <p>
          Мы стремимся дать не только узкие знания и умения, но и научить общаться, проявлять себя,
          ставить и добиваться целей. Мы называем себя дизайн-центром: дизайн — это художественное
          конструирование предметов, и именно этим мы с ребятами и занимаемся.
        </p>
      </div>
      <PhotoFrame ratio="16 / 9" src="/media/photos/090.jpg" alt="Мастер-класс в центре «Мы Друзья»" />
    </div>
  </section>

  <!-- Принципы -->
  <section class="section section--muted">
    <div class="container">
      <SectionHeader eyebrow="Наши принципы" title="Во что мы верим" center />
      <div class="grid" style="--min: 260px">
        <article v-for="p in principles" :key="p.title" class="principle">
          <span class="principle__emoji" aria-hidden="true">{{ p.emoji }}</span>
          <h3 class="principle__title">{{ p.title }}</h3>
          <p class="principle__text">{{ p.text }}</p>
        </article>
      </div>
    </div>
  </section>

  <!-- Коротко о нас -->
  <section class="section">
    <div class="container container--narrow">
      <SectionHeader eyebrow="Коротко" title="Чем мы живём" />
      <ul class="facts">
        <li v-for="fact in facts" :key="fact.text" class="facts__item">
          <span class="facts__emoji" aria-hidden="true">{{ fact.emoji }}</span>
          <p class="facts__text">{{ fact.text }}</p>
        </li>
      </ul>
    </div>
  </section>

  <CtaSection />
</template>

<style scoped>
.principle {
  padding: var(--space-6);
  text-align: center;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.principle__emoji {
  font-size: 40px;
}
.principle__title {
  font-size: var(--text-lg);
  margin-block: var(--space-3) var(--space-2);
}
.principle__text {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}

/* Факты о центре */
.facts {
  display: grid;
  gap: var(--space-3);
}
.facts__item {
  display: flex;
  gap: var(--space-4);
  align-items: center;
  padding: var(--space-4) var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-left: 4px solid var(--color-accent);
  border-radius: var(--radius-md);
}
.facts__emoji {
  font-size: 28px;
  flex-shrink: 0;
}
.facts__text {
  color: var(--color-text-muted);
}
</style>
