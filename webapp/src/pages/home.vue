<script setup>

import {ref} from "vue";
import {useFormattedDate} from '../hooks/useFormattedDate';

const {getDate} = useFormattedDate();
const services = ref([])
const servicesLoaded = ref(false)

async function getServicesFetch() {
  let url = `${import.meta.env.VITE_API_BASE}/me/services`;
  let method = 'GET';

  try {
    const response = await fetch(url, {
      method: method,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    switch (response.status) {
      case 401:
        console.error("Не авторизован");
        return;
      case 200:
        const text = await response.text();
        if (text) {
          services.value = JSON.parse(text).reverse();
        }
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }

  servicesLoaded.value = true;
}

getServicesFetch()
</script>

<template>
  <div>
    <div class="head flex items-center justify-between">
      <h1 class="text-3xl font-extrabold">Мои услуги
      </h1>
      <router-link :to="{name:'create'}" class="bg-black text-white p-2 rounded-md">➕ Заказать</router-link>
    </div>
    <div class="content mt-3">
      <div class="grid grid-cols-1 gap-4">
        <div v-for="item in [1,2,3]" :class="!services?'':'skeleton'" v-if="!services[0] && !servicesLoaded"
             class="p-6">
          Lorem ipsum dolor sit amet, consectetur adipisicing elit. Aspernatur est eum ipsum placeat porro sapiente
          soluta temporibus. Beatae dolores ducimus fuga sequi voluptas voluptatibus! Aut culpa doloremque eius in
          ullam?
        </div>
        <router-link
            :to="`/services?id=${service.id}`"
            class="bg-white rounded-xl shadow-lg p-6 hover:shadow-xl transition-shadow duration-300 ease-in-out cursor-pointer"
            v-for="service in services"
            :key="service.id"
        >
          <h3 class="text-lg font-semibold text-gray-800">{{ service.name }}</h3>
          <p class="text-gray-600 mt-1">{{ service.location }}</p>
          <p class="text-gray-600 mt-1">Оплачено до: <span class="font-bold">{{ getDate(service.expires_at) }}</span>
          </p>
          <p class="text-gray-600 mt-1">Сервис: <span class="font-bold">{{ service.service }}</span></p>
        </router-link>
      </div>

    </div>
  </div>
</template>

<style scoped>

</style>
