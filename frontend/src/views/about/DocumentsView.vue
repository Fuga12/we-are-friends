<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import { documents } from '@/data/documents'
</script>

<template>
  <PageHero
    title="Документы"
    subtitle="Лицензии, устав и другие документы центра. Всё открыто и доступно для ознакомления."
    emoji="📄"
  />

  <section class="section">
    <div class="container container--narrow">
      <ul class="doc-list">
        <li v-for="doc in documents" :key="doc.id">
          <a
            class="doc"
            :href="doc.url"
            target="_blank"
            rel="noopener noreferrer"
            :aria-label="`Скачать: ${doc.title}`"
          >
            <span class="doc__icon" aria-hidden="true">📄</span>
            <span class="doc__body">
              <span class="doc__title">{{ doc.title }}</span>
              <span class="doc__desc">{{ doc.description }}</span>
            </span>
            <span class="doc__meta">
              <span class="doc__type">{{ doc.fileType }}</span>
              <span class="doc__size">{{ doc.fileSize }}</span>
            </span>
            <span class="doc__download" aria-hidden="true">↓</span>
          </a>
        </li>
      </ul>

      <p class="text-muted text-center doc-note">
        Ссылки — заглушки. Загрузи реальные файлы и подставь адреса в
        <code>src/data/documents.ts</code>.
      </p>
    </div>
  </section>
</template>

<style scoped>
.doc-list {
  display: grid;
  gap: var(--space-3);
}
.doc {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4) var(--space-5);
  background-color: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  transition: border-color var(--transition), box-shadow var(--transition), transform var(--transition);
}
.doc:hover {
  color: var(--color-text);
  border-color: var(--color-primary);
  box-shadow: var(--shadow-sm);
  transform: translateX(3px);
}
.doc__icon {
  font-size: 28px;
  flex-shrink: 0;
}
.doc__body {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}
.doc__title {
  font-family: var(--font-heading);
  font-weight: var(--weight-semibold);
}
.doc__desc {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}
.doc__meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  flex-shrink: 0;
  font-size: var(--text-xs);
  color: var(--color-text-muted);
}
.doc__type {
  font-weight: var(--weight-bold);
  color: var(--color-accent-strong);
}
.doc__download {
  display: grid;
  place-items: center;
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  color: var(--color-primary);
  background-color: var(--color-primary-soft);
  border-radius: 50%;
  font-weight: var(--weight-bold);
}
.doc-note {
  margin-top: var(--space-6);
  font-size: var(--text-sm);
}
.doc-note code {
  font-size: 0.9em;
  padding: 1px 6px;
  background-color: var(--color-surface-muted);
  border-radius: var(--radius-sm);
}

/* На мобилке прячем колонку с типом/размером — оставляем главное */
@media (max-width: 520px) {
  .doc__meta {
    display: none;
  }
}
</style>
