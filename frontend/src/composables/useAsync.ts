/* Маленький помощник для загрузки данных с бэкенда.

   Избавляет компоненты от повторяющегося кода: сам держит три состояния —
   data (что пришло), loading (идёт загрузка) и error (текст ошибки).

   Пример в компоненте:
     const { data: news, loading, error } = useAsync(getNews)
   Пока данные едут — loading === true; когда пришли — они в news.value.

   Второй аргумент watchSource нужен для страниц, зависящих от адреса
   (например /news/:id): передайте () => route.params.id — и при переходе
   на другой id данные перезагрузятся сами. */
import { ref, shallowRef, watch, type WatchSource } from 'vue'

export function useAsync<T>(loader: () => Promise<T>, watchSource?: WatchSource) {
  const data = shallowRef<T | null>(null)
  const loading = ref(true)
  const error = ref<string | null>(null)

  async function run() {
    loading.value = true
    error.value = null
    try {
      data.value = await loader()
    } catch (e) {
      data.value = null
      error.value = e instanceof Error ? e.message : 'Не удалось загрузить данные'
    } finally {
      loading.value = false
    }
  }

  // Перезагружаем при смене ключа (например id в адресе).
  if (watchSource) watch(watchSource, run)

  run()

  return { data, loading, error, reload: run }
}
