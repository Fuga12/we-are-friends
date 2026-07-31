<script setup lang="ts">
import { mainNav } from '@/data/navigation'
import { site } from '@/data/site'
import { Clock3, MapPin, Phone } from 'lucide-vue-next'

const year = new Date().getFullYear()
</script>

<template>
  <footer class="footer">
    <div class="container footer__grid">
      <!-- О центре -->
      <div class="footer__col footer__about">
        <span class="footer__logo">{{ site.name }}</span>
        <p class="footer__desc">{{ site.description }}</p>
        <a class="footer__social" :href="site.vk" target="_blank" rel="noopener noreferrer">
          Мы во ВКонтакте →
        </a>
      </div>

      <!-- Навигация -->
      <nav class="footer__col" aria-label="Меню в подвале">
        <h3 class="footer__heading">Разделы</h3>
        <ul class="footer__links">
          <li v-for="link in mainNav" :key="link.to">
            <RouterLink :to="link.to">{{ link.label }}</RouterLink>
          </li>
        </ul>
      </nav>

      <!-- Контакты -->
      <div class="footer__col">
        <h3 class="footer__heading">Контакты</h3>
        <ul class="footer__contacts">
          <li>
            <MapPin :size="18" stroke-width="1.8" aria-hidden="true" />
            <a :href="site.yandexMaps" target="_blank" rel="noopener noreferrer">
              {{ site.address }}
            </a>
          </li>
          <li>
            <Phone :size="18" stroke-width="1.8" aria-hidden="true" />
            <a :href="site.phoneHref">{{ site.phone }}</a>
          </li>
          <li>
            <Phone :size="18" stroke-width="1.8" aria-hidden="true" />
            <a :href="site.phone2Href">{{ site.phone2 }}</a>
          </li>
          <li>
            <Clock3 :size="18" stroke-width="1.8" aria-hidden="true" />
            <span>{{ site.workingHours }}</span>
          </li>
        </ul>
      </div>
    </div>

    <div class="footer__bottom">
      <div class="container footer__bottom-inner">
        <span>© {{ year }} Детский центр «{{ site.name }}»</span>
        <RouterLink to="/about/documents">Документы и политика</RouterLink>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  margin-top: var(--section-gap);
  background-color: var(--brand-teal);
  color: var(--color-text-inverse);
}
.footer__grid {
  display: grid;
  gap: var(--space-6);
  padding-block: var(--space-8) var(--space-6);
  grid-template-columns: repeat(auto-fit, minmax(min(220px, 100%), 1fr));
}
.footer__about {
  max-width: 34ch;
}
.footer__logo {
  font-family: var(--font-heading);
  font-size: var(--text-2xl);
  font-weight: var(--weight-bold);
}
.footer__desc {
  margin-top: var(--space-3);
  color: rgba(255, 255, 255, 0.8);
  font-size: var(--text-sm);
}
.footer__social {
  display: inline-block;
  margin-top: var(--space-4);
  font-weight: var(--weight-bold);
  color: var(--brand-amber);
}
.footer__social:hover {
  color: #fff;
}
.footer__heading {
  font-size: var(--text-lg);
  margin-bottom: var(--space-4);
  color: #fff;
}
.footer__links a,
.footer__contacts a {
  color: rgba(255, 255, 255, 0.85);
}
.footer__links a:hover,
.footer__contacts a:hover {
  color: var(--brand-amber);
}
.footer__links li + li {
  margin-top: var(--space-2);
}
.footer__contacts {
  display: grid;
  gap: var(--space-3);
}
.footer__contacts li {
  display: flex;
  gap: var(--space-3);
  color: rgba(255, 255, 255, 0.85);
  font-size: var(--text-sm);
}
.footer__contacts svg {
  flex: none;
  color: var(--brand-amber);
}
.footer__bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.15);
}
.footer__bottom-inner {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
  justify-content: space-between;
  padding-block: var(--space-4);
  font-size: var(--text-sm);
  color: rgba(255, 255, 255, 0.7);
}
.footer__bottom-inner a {
  color: rgba(255, 255, 255, 0.7);
}
.footer__bottom-inner a:hover {
  color: #fff;
}
</style>
