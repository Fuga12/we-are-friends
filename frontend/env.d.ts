/// <reference types="vite/client" />

// Позволяет TypeScript понимать импорты .vue файлов
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
