<template>
  <div class="fixed right-3 2xl:right-auto flex flex-col items-center justify-center z-50 min-w-64">
    <div
        v-for="(toast, index) in toasts"
        :key="index"
        :class="[
        'm-3 p-4 rounded-md shadow-lg text-white w-full',
        toast.severityClass,
        toast.visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
      ]"
        style="transition: all 0.3s ease;"
    >
      <h3 class="font-bold text-lg mb-2">{{ toast.summary }}</h3>
      <p>{{ toast.detail }}</p>
    </div>
  </div>
</template>

<script setup>
import {useToast} from '../hooks/useToast';

// Получаем массив тостов
const {toasts} = useToast();
</script>
