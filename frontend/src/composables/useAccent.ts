import { computed, type MaybeRefOrGetter, toValue } from 'vue'

/*
  Многие карточки бывают четырёх фирменных цветов (teal/orange/amber/red).
  Вместо того чтобы плодить CSS-классы, отдаём объект со стилями,
  где заданы CSS-переменные --c (основной) и --c-soft (подложка).
  В стилях компонента используем var(--c) / var(--c-soft).
*/
export type AccentColor = 'teal' | 'orange' | 'amber' | 'red'

const map: Record<AccentColor, { c: string; soft: string; strong: string }> = {
  teal: { c: 'var(--brand-teal)', soft: 'var(--brand-teal-soft)', strong: 'var(--brand-teal-strong)' },
  orange: { c: 'var(--brand-orange)', soft: 'var(--brand-orange-soft)', strong: 'var(--brand-orange-strong)' },
  amber: { c: 'var(--brand-amber)', soft: 'var(--brand-amber-soft)', strong: 'var(--brand-red-strong)' },
  red: { c: 'var(--brand-red)', soft: 'var(--brand-red-soft)', strong: 'var(--brand-red-strong)' },
}

export function useAccent(color: MaybeRefOrGetter<AccentColor>) {
  return computed(() => {
    const value = map[toValue(color)] ?? map.teal
    return {
      '--c': value.c,
      '--c-soft': value.soft,
      '--c-strong': value.strong,
    } as Record<string, string>
  })
}
