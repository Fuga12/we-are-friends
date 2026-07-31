import { createApp } from 'vue'
import Aura from '@primeuix/themes/aura'
import PrimeVue from 'primevue/config'
import App from './App.vue'
import router from './router'
import './assets/styles/main.css'

createApp(App)
  .use(PrimeVue, { theme: { preset: Aura, options: { darkModeSelector: false } } })
  .use(router)
  .mount('#app')
