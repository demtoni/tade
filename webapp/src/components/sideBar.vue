<template>
  <!-- Затемняющий слой -->
  <div
      v-if="enabled"
      class="fixed inset-0 bg-black bg-opacity-50 z-40 transition-opacity duration-300"
      @click="closePanel"
  ></div>

  <!-- Панель -->
  <div
      v-if="enabled"
      class="fixed top-0 h-full w-64 bg-white shadow-lg transform transition-transform z-50"
      :class="panelClasses"
      @transitionend="handleTransitionEnd"
  >
    <div class="flex justify-between items-center p-4 border-b">
      <h2 class="text-lg font-bold">{{ title }}</h2>
      <button @click="closePanel" class="text-xl font-bold">&times;</button>
    </div>
    <div class="p-4">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'

const props = defineProps({
  enabled: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Панель'
  },
  position: {
    type: String,
    default: 'right',
    validator: (value) => ['left', 'right'].includes(value)
  }
})

const emit = defineEmits(['update:enabled'])

function closePanel() {
  emit('update:enabled', false)
}

function handleTransitionEnd() {
  if (!props.enabled) {
    emit('update:enabled', false)
  }
}

const panelClasses = computed(() => {
  const baseClass = !props.enabled ? 'duration-200' : 'duration-500'
  const positionClass = props.position === 'right' ? 'right-0' : 'left-0'
  const translateClass = props.position === 'right'
      ? (!props.enabled ? 'translate-x-full' : 'translate-x-0')
      : (!props.enabled ? '-translate-x-full' : 'translate-x-0')

  return `${baseClass} ${positionClass} ${translateClass}`
})

watch(() => props.enabled, (newVal) => {
  if (!newVal) {
    emit('update:enabled', false)
  }
})
</script>

<style scoped>
/* Дополнительные стили (если нужны) */
</style>
