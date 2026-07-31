<script setup lang="ts">
import { reactive, ref } from 'vue'
import { Clock3, House, MapPin, MessageCircle, Phone } from 'lucide-vue-next'
import PageHero from '@/components/layout/PageHero.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import { site } from '@/data/site'

/* Простая форма записи. Бэкенда пока нет — по отправке показываем
   сообщение об успехе. Когда появится API, замени тело submit() на fetch(). */
const form = reactive({ name: '', phone: '', message: '' })
const errors = reactive<{ name?: string; phone?: string }>({})
const sent = ref(false)

function validate() {
  errors.name = form.name.trim() ? undefined : 'Пожалуйста, укажите имя'
  errors.phone = /[\d]{5,}/.test(form.phone) ? undefined : 'Укажите корректный телефон'
  return !errors.name && !errors.phone
}

function submit() {
  if (!validate()) return
  // TODO: отправить на бэкенд, например:
  // await fetch('/api/leads', { method: 'POST', body: JSON.stringify(form) })
  sent.value = true
}

const contactCards = [
  { icon: MapPin, color: 'var(--brand-red)', label: 'Основная площадка', value: site.address },
  { icon: House, color: 'var(--brand-teal)', label: 'Филиал (Древлянка)', value: `${site.address2} — ${site.address2Note}` },
  { icon: Phone, color: 'var(--brand-amber)', label: 'Телефон 1', value: site.phone, href: site.phoneHref },
  { icon: Phone, color: 'var(--brand-amber)', label: 'Телефон 2', value: site.phone2, href: site.phone2Href },
  { icon: MessageCircle, color: 'var(--brand-blue)', label: 'ВКонтакте', value: 'vk.ru/ddc_ptz', href: site.vk },
  { icon: Clock3, color: 'var(--brand-orange)', label: 'Часы работы', value: site.workingHours },
]
</script>

<template>
  <PageHero
    title="Контакты"
    subtitle="Позвоните, напишите или приходите в гости — мы всегда рады знакомству."
    emoji="📍"
  />

  <section class="section">
    <div class="container contacts">
      <!-- Левая колонка: контакты + карта -->
      <div class="contacts__info stack-lg">
        <ul class="contacts__cards">
          <li v-for="c in contactCards" :key="c.label" class="contact-card">
            <component
              :is="c.icon"
              class="contact-card__icon"
              :style="{ color: c.color }"
              :size="26"
              stroke-width="1.8"
              aria-hidden="true"
            />
            <div>
              <span class="contact-card__label">{{ c.label }}</span>
              <a v-if="c.href" :href="c.href" class="contact-card__value">{{ c.value }}</a>
              <span v-else class="contact-card__value">{{ c.value }}</span>
            </div>
          </li>
        </ul>

        <!-- Карта проезда: виджет Яндекс.Карт.
             Размеры не задаём атрибутами width/height — их держит CSS (.map),
             иначе фиксированная ширина iframe растягивает колонку сетки. -->
        <div class="map">
          <iframe
            src="https://yandex.ru/map-widget/v1/?um=constructor%3A80522a5598c600b03ad678a76cbd16d89bf72f00c63e1228c017bef153e91ff5&amp;source=constructor"
            title="Карта проезда до центра «Мы Друзья»"
            loading="lazy"
            allowfullscreen
          ></iframe>
        </div>
      </div>

      <!-- Правая колонка: форма -->
      <div class="contacts__form-wrap">
        <div v-if="!sent" class="form-card">
          <h2 class="form-card__title">Записаться на занятие</h2>
          <p class="text-muted">Оставьте контакты — мы перезвоним и всё расскажем.</p>

          <form class="form" novalidate @submit.prevent="submit">
            <div class="field">
              <label for="name">Имя</label>
              <input
                id="name"
                v-model="form.name"
                type="text"
                placeholder="Как вас зовут?"
                :class="{ 'field--error': errors.name }"
                :aria-invalid="!!errors.name"
              />
              <span v-if="errors.name" class="field__error">{{ errors.name }}</span>
            </div>

            <div class="field">
              <label for="phone">Телефон</label>
              <input
                id="phone"
                v-model="form.phone"
                type="tel"
                placeholder="+7 (___) ___-__-__"
                :class="{ 'field--error': errors.phone }"
                :aria-invalid="!!errors.phone"
              />
              <span v-if="errors.phone" class="field__error">{{ errors.phone }}</span>
            </div>

            <div class="field">
              <label for="message">Сообщение <span class="text-muted">(необязательно)</span></label>
              <textarea
                id="message"
                v-model="form.message"
                rows="4"
                placeholder="Какое направление интересует, возраст ребёнка…"
              ></textarea>
            </div>

            <BaseButton type="submit" variant="accent" size="lg" block>Отправить заявку</BaseButton>
            <p class="form__note text-muted">
              Нажимая кнопку, вы соглашаетесь с обработкой персональных данных.
            </p>
          </form>
        </div>

        <!-- Экран успеха -->
        <div v-else class="form-card form-card--success">
          <span class="form-card__check" aria-hidden="true">✅</span>
          <h2 class="form-card__title">Спасибо, {{ form.name }}!</h2>
          <p>Мы получили вашу заявку и скоро перезвоним по номеру {{ form.phone }}.</p>
          <BaseButton to="/" variant="outline">На главную</BaseButton>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.contacts {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-7);
  align-items: start;
}
/* min-width: 0 разрешает колонкам сжиматься. Без этого широкий вложенный
   элемент (например iframe карты) раздвигает колонку и ломает сетку. */
.contacts > * {
  min-width: 0;
}

/* Контактные карточки */
.contacts__cards {
  display: grid;
  gap: var(--space-3);
}
.contact-card {
  display: flex;
  gap: var(--space-4);
  padding: var(--space-4);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}
.contact-card__icon {
  flex: none;
}
.contact-card__label {
  display: block;
  font-size: var(--text-xs);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--color-text-muted);
}
.contact-card__value {
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  color: var(--color-text);
}
a.contact-card__value:hover {
  color: var(--color-primary);
}

/* Карта проезда (виджет Яндекс.Карт).
   Ширина — по колонке, высота — пропорцией, поэтому карта одинаково хорошо
   ложится и в две колонки на десктопе, и в одну на мобилке. */
.map {
  width: 100%;
  aspect-ratio: 3 / 2;
  overflow: hidden;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  /* подложка видна, пока виджет подгружается */
  background-color: var(--color-surface-muted);
}
.map iframe {
  display: block;
  width: 100%;
  height: 100%;
  border: 0;
  /* Firefox не всегда обрезает поверхность iframe через overflow родителя.
     Клиппинг на самом iframe убирает торчащие углы карты. */
  clip-path: inset(0 round var(--radius-lg));
}

/* Форма */
.form-card {
  padding: var(--space-6);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-md);
}
.form-card__title {
  font-size: var(--text-2xl);
  margin-bottom: var(--space-2);
}
.form {
  display: grid;
  gap: var(--space-4);
  margin-top: var(--space-5);
}
.field {
  display: grid;
  gap: var(--space-2);
}
.field label {
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
  font-size: var(--text-sm);
}
.field input,
.field textarea {
  width: 100%;
  padding: var(--space-3) var(--space-4);
  background-color: var(--color-bg);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: border-color var(--transition), box-shadow var(--transition);
}
.field input:focus,
.field textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: var(--shadow-focus);
}
.field textarea {
  resize: vertical;
}
.field--error {
  border-color: var(--color-error) !important;
}
.field__error {
  font-size: var(--text-sm);
  color: var(--color-error);
}
.form__note {
  font-size: var(--text-xs);
  text-align: center;
}

/* Экран успеха */
.form-card--success {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--space-4);
  padding-block: var(--space-8);
}
.form-card__check {
  font-size: 56px;
}

/* Адаптив: одна колонка на планшете и мобилке */
@media (max-width: 860px) {
  .contacts {
    grid-template-columns: 1fr;
  }
}
</style>
