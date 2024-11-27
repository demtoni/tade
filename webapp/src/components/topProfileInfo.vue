<script setup>
import {useUser} from '../hooks/useUser.js';
import {onMounted} from 'vue';

const {meData, getMeFetch, logout} = useUser();

onMounted(() => {
  getMeFetch();
});

</script>

<template>
  <div class="flex gap-5 items-center">
    <p class="bg-gray-100 rounded-xl py-2 px-6 text-center" :class="meData && meData.username?'':'skeleton'"><span v-if="meData" >{{ meData.username }}</span>
      <span v-else>Загрузка...</span></p>
    <p class="cursor-pointer select-none font-semibold hover:text-blue-900" @click="logout()">Выйти</p>
  </div>
</template>

<style scoped>

</style>