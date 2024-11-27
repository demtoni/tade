<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center" v-if="isOpen" @click="closeModal">
    <div class="bg-white rounded-lg shadow-lg w-96" @click.stop>
      <div class="flex justify-between items-center p-4 border-b">
        <h2 class="text-xl font-semibold">{{ title }}</h2>
        <button class="text-gray-500 hover:text-gray-800" @click="closeModal">&times;</button>
      </div>
      <div class="p-4">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  isOpen: {
    type: Boolean,
    required: true,
  },
});

// Определяем события
const emit = defineEmits();

// Функция для закрытия модального окна
const closeModal = () => {
  emit('update:isOpen', false);
};
</script>

<style scoped>
/* Дополнительные стили могут быть добавлены здесь */
</style>
