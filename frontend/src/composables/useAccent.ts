import { computed, type MaybeRefOrGetter, toValue } from 'vue'

/*
  Акцентный цвет карточки приходит с бэкенда как HEX-строка ('#1a6757').

  Компоненту нужны три оттенка, поэтому здесь мы отдаём объект с CSS-переменными:
    --c        — сам цвет как есть;
    --c-soft   — светлая подложка (цвет, разбавленный фоном страницы);
    --c-strong — затемнённый оттенок для текста и hover.

  Оттенки считает сам браузер через color-mix() — так бэкенду достаточно
  прислать один HEX, и работает любой цвет, а не только четыре фирменных.
  В стилях компонента используем var(--c) / var(--c-soft) / var(--c-strong).
*/
const FALLBACK = 'var(--brand-teal)'

export function useAccent(color: MaybeRefOrGetter<string | undefined>) {
  return computed(() => {
    const c = toValue(color) || FALLBACK
    return {
      '--c': c,
      // 12% цвета на фоне страницы — та же лёгкая подложка, что была у brand-*-soft
      '--c-soft': `color-mix(in srgb, ${c} 12%, var(--color-bg))`,
      // ~18% чёрного — примерно такой же шаг, как между brand-teal и brand-teal-strong
      '--c-strong': `color-mix(in srgb, ${c} 82%, #000)`,
    } as Record<string, string>
  })
}
