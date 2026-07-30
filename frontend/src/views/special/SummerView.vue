<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import PhotoFrame from '@/components/ui/PhotoFrame.vue'
import CtaSection from '@/components/sections/CtaSection.vue'

const day = [
  { time: '9:00', text: 'Сбор, зарядка и планы на день' },
  { time: '10:00', text: 'Творческий блок или мастер-класс' },
  { time: '12:00', text: 'Прогулка и активные игры на улице' },
  { time: '13:00', text: 'Обед' },
  { time: '14:00', text: 'Эксперименты, квесты, настолки' },
  { time: '16:00', text: 'Полдник и свободные игры' },
  { time: '17:00', text: 'Итоги дня, встречаем родителей' },
]

/* Реальные летние программы центра (vk.ru/ddc_ptz). */
const shifts = [
  {
    emoji: '🪄',
    title: 'Школа магии: Тайная миссия',
    dates: '13–24 июля',
    photo: '/media/photos/177.jpg',
    text: 'Ребята становятся учениками школы волшебства: захватывающий квест, удивительные опыты и магические изделия на столярном мастер-классе. А ещё каждый выращивает собственную микрозелень и забирает садик домой.',
    activities: ['Квест и магические испытания', 'Опыты и эксперименты', 'Столярный мастер-класс', 'Своя микрозелень'],
  },
  {
    emoji: '💖',
    title: 'Деффчонки',
    dates: '27 июля – 7 августа',
    photo: '/media/photos/160.jpg',
    text: 'Самая нежная смена: творческие мастер-классы, изделия своими руками и профессиональная фотосессия (400+ снимков). Уход за собой, стиль и уверенность, а также поход, боулинг, лазертаг и пижамная вечеринка.',
    activities: ['Профессиональная фотосессия', 'Творческие мастер-классы', 'Стиль и уход за собой', 'Поход, боулинг, лазертаг'],
  },
  {
    emoji: '🤖',
    title: 'Майнкрафт: от игры к роботу',
    dates: '10–28 августа',
    photo: '/media/photos/050.jpg',
    text: 'Для любителей Minecraft: знакомство с робототехникой и программированием, создание мини-роботов и гонки роботов «Формула 1». Плюс бассейн, походы, лазертаг и кино. Можно выбрать участие на 1, 2 или 3 недели.',
    activities: ['Робототехника и программирование', 'Гонки роботов «Формула 1»', 'Бассейн и походы', 'Участие на 1–3 недели'],
  },
]
</script>

<template>
  <PageHero
    title="Летние каникулы"
    subtitle="Городской лагерь дневного пребывания: ребёнок занят и счастлив, а родители спокойны."
    emoji="☀️"
    tone="amber"
  />

  <!-- Смены -->
  <section class="section">
    <div class="container">
      <SectionHeader eyebrow="Смены" title="Программы на лето" center />
      <div class="grid" style="--min: 300px">
        <article v-for="s in shifts" :key="s.title" class="shift">
          <div class="shift__media">
            <PhotoFrame ratio="16 / 10" :src="s.photo" :alt="s.title" :emoji="s.emoji" />
            <span class="shift__dates">{{ s.dates }}</span>
          </div>
          <div class="shift__body">
            <h3 class="shift__title">{{ s.title }}</h3>
            <p class="shift__text">{{ s.text }}</p>
            <ul class="shift__list">
              <li v-for="a in s.activities" :key="a">{{ a }}</li>
            </ul>
          </div>
        </article>
      </div>
    </div>
  </section>

  <!-- Распорядок дня -->
  <section class="section section--muted">
    <div class="container container--narrow">
      <SectionHeader eyebrow="Как проходит день" title="Примерный распорядок" />
      <ul class="daily">
        <li v-for="item in day" :key="item.time" class="daily__row">
          <span class="daily__time">{{ item.time }}</span>
          <span class="daily__text">{{ item.text }}</span>
        </li>
      </ul>
    </div>
  </section>

  <CtaSection
    title="Забронируйте место в смене"
    text="Группы небольшие и заполняются быстро. Позвоните, чтобы занять место для ребёнка."
  />
</template>

<style scoped>
.shift {
  display: flex;
  flex-direction: column;
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-top: 5px solid var(--brand-amber);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  height: 100%;
}
.shift__media {
  position: relative;
}
.shift__media :deep(.photo) {
  border-radius: 0;
  --c-soft: var(--color-highlight-soft); /* тёплая летняя подложка заглушки */
}
.shift__dates {
  position: absolute;
  top: var(--space-3);
  left: var(--space-3);
  padding: 3px 12px;
  font-size: var(--text-sm);
  font-weight: var(--weight-semibold);
  color: var(--brand-red-strong);
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: var(--radius-pill);
  box-shadow: var(--shadow-sm);
}
.shift__body {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-5);
}
.shift__title {
  font-size: var(--text-xl);
}
.shift__text {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}
.shift__list {
  display: grid;
  gap: var(--space-2);
  padding-top: var(--space-3);
  border-top: 1px dashed var(--color-border);
}
.shift__list li {
  position: relative;
  padding-left: 1.6em;
  font-size: var(--text-sm);
}
.shift__list li::before {
  content: '✓';
  position: absolute;
  left: 0;
  font-weight: var(--weight-bold);
  color: var(--brand-amber);
}

.daily {
  display: grid;
  gap: var(--space-1);
}
.daily__row {
  display: flex;
  gap: var(--space-4);
  align-items: baseline;
  padding: var(--space-3) var(--space-4);
  background-color: var(--color-surface);
  border-radius: var(--radius-md);
}
.daily__time {
  flex-shrink: 0;
  width: 4.5em;
  font-family: var(--font-heading);
  font-weight: var(--weight-bold);
  color: var(--color-primary);
}
.daily__text {
  color: var(--color-text-muted);
}
</style>
