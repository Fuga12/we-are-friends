<script setup lang="ts">
import SiteHeader from '@/components/layout/SiteHeader.vue'
import SiteFooter from '@/components/layout/SiteFooter.vue'
</script>

<template>
  <SiteHeader />

  <main id="main">
    <!-- Плавное перекрёстное затухание при смене страницы.
         :key заставляет Vue считать каждую страницу новой и проигрывать анимацию.
         Здесь НЕТ mode="out-in": новая страница появляется одновременно с уходом
         старой, поэтому пустого («белого») экрана между ними не бывает. -->
    <RouterView v-slot="{ Component, route }">
      <Transition name="page">
        <component :is="Component" :key="route.path" />
      </Transition>
    </RouterView>
  </main>

  <SiteFooter />
</template>

<style scoped>
#main {
  position: relative; /* точка отсчёта для уходящей страницы (см. ниже) */
  min-height: 60vh;
}
.page-enter-active,
.page-leave-active {
  transition: opacity 0.28s var(--ease);
}
/* Уходящую страницу вынимаем из потока и кладём поверх — тогда новая
   сразу занимает место, а старая плавно растворяется над ней. */
.page-leave-active {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
}
.page-enter-from,
.page-leave-to {
  opacity: 0;
}
</style>
