<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import CtaSection from '@/components/sections/CtaSection.vue'

const topics = [
  { emoji: '🎨', title: 'Роспись', text: 'Раскрашиваем гипсовые фигурки, футболки или деревянные игрушки.' },
  { emoji: '🧼', title: 'Мыловарение', text: 'Делаем ароматное мыло разных форм и цветов своими руками.' },
  { emoji: '🏺', title: 'Лепка из глины', text: 'Керамика и пластилин: развиваем мелкую моторику и фантазию.' },
  { emoji: '🕯️', title: 'Свечи', text: 'Создаём декоративные свечи — отличный подарок близким.' },
  { emoji: '🌸', title: 'Флористика', text: 'Собираем композиции и открытки из сухоцветов.' },
  { emoji: '🔬', title: 'Научное шоу', text: 'Безопасные и зрелищные опыты, которые удивят детей.' },
]
</script>

<template>
  <PageHero
    title="Мастер-классы"
    subtitle="Творческие занятия у нас в центре и выездные — в школах, детских садах и на праздниках."
    emoji="🎭"
  />

  <!-- Формат -->
  <section class="section">
    <div class="container container--narrow stack">
      <p class="lead">
        Каждый мастер-класс — это готовый маленький праздник творчества. Мы привозим все материалы,
        а каждый ребёнок уносит домой поделку, сделанную своими руками.
      </p>
      <div class="formats">
        <div class="format">
          <span class="format__emoji" aria-hidden="true">🏠</span>
          <div>
            <h3>У нас в центре</h3>
            <p class="text-muted">Приходите группой или запишитесь на открытый мастер-класс.</p>
          </div>
        </div>
        <div class="format">
          <span class="format__emoji" aria-hidden="true">🚌</span>
          <div>
            <h3>Выездной формат</h3>
            <p class="text-muted">Приедем к вам: в школу, сад или на день рождения. От 5 участников.</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Темы -->
  <section class="section section--muted">
    <div class="container">
      <SectionHeader eyebrow="Направления" title="Что можно попробовать" center />
      <div class="grid" style="--min: 240px">
        <article v-for="t in topics" :key="t.title" class="topic">
          <span class="topic__emoji" aria-hidden="true">{{ t.emoji }}</span>
          <h3 class="topic__title">{{ t.title }}</h3>
          <p class="topic__text">{{ t.text }}</p>
        </article>
      </div>
    </div>
  </section>

  <CtaSection
    title="Закажите мастер-класс"
    text="Расскажите, для какого возраста и сколько детей — подберём тему и рассчитаем стоимость."
  />
</template>

<style scoped>
.formats {
  display: grid;
  gap: var(--space-4);
  grid-template-columns: repeat(auto-fit, minmax(min(260px, 100%), 1fr));
}
.format {
  display: flex;
  gap: var(--space-4);
  padding: var(--space-5);
  background-color: var(--color-primary-soft);
  border-radius: var(--radius-lg);
}
.format__emoji {
  font-size: 34px;
}
.format h3 {
  font-size: var(--text-lg);
  margin-bottom: var(--space-1);
}

.topic {
  padding: var(--space-5);
  text-align: center;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}
.topic__emoji {
  font-size: 40px;
}
.topic__title {
  font-size: var(--text-lg);
  margin-block: var(--space-2);
}
.topic__text {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
</style>
