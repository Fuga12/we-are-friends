<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import Button from 'primevue/button'

withDefaults(defineProps<{ label?: string; loading?: boolean; size?: 'small' | 'large' }>(), {
  label: 'Удалить',
  loading: false,
})

const emit = defineEmits<{ confirm: [] }>()
const armed = ref(false)
let resetTimer: ReturnType<typeof setTimeout> | undefined

function click() {
  if (!armed.value) {
    armed.value = true
    resetTimer = setTimeout(() => { armed.value = false }, 5000)
    return
  }
  if (resetTimer) clearTimeout(resetTimer)
  armed.value = false
  emit('confirm')
}

onBeforeUnmount(() => { if (resetTimer) clearTimeout(resetTimer) })
</script>

<template>
  <Button :label="armed ? 'Точно удалить?' : label" :loading="loading" :size="size" severity="danger" text @click="click" />
</template>
