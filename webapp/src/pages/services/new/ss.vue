<script setup>
import {ref, inject, computed, onMounted, watch} from "vue";
import {useRouter} from 'vue-router';
import {useUser} from "../../../hooks/useUser.js";
import Modal from "../../../components/modal.vue";

const router = useRouter();
const addToast = inject('addToast');
const {meData, getMeFetch} = useUser();


onMounted(() => {
  getMeFetch();
});

const locations = ref([]);
const algorithms = ref(['chacha20-ietf-poly1305']);

const form = ref({
  name: generateName(),
  location: 'none',
  selectedMonth: 1,
  selectedPlugin: 'none',
  selectedAlgorithm: 'chacha20-ietf-poly1305',
  prolong: false
});

const total = computed(() => form.value.selectedMonth * 100);

function handleSubmit() {
  if (meData.value.balance < total.value) {
    openModal()
    return
  }
  createService();
}

function generateName() {
  const adjectives = [
    'pink', 'regular', 'happy', 'strong', 'fast', 'silent', 'brave', 'bright', 'wild', 'calm'
  ];

  const nouns = [
    'boat', 'memory', 'cloud', 'star', 'river', 'forest', 'mountain', 'stone', 'wind', 'fire'
  ];

  const randomAdjective = adjectives[Math.floor(Math.random() * adjectives.length)];
  const randomNoun = nouns[Math.floor(Math.random() * nouns.length)];
  const randomNumber = Math.floor(Math.random() * 10);

  return `${randomAdjective}-${randomNoun}-${randomNumber}`;
}

async function getLocationsFetch() {
  let url = `${import.meta.env.VITE_API_BASE}/me/services/locations`;
  let method = 'GET';

  try {
    const response = await fetch(url, {
      method: method,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      }
    });

    switch (response.status) {
      case 401:
        console.error("Не авторизован");
        return;
      case 200:
        const text = await response.text();
        if (text) {
          locations.value = JSON.parse(text);
        }
        break;
      default:
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }
}

async function createService() {
  let url = `${import.meta.env.VITE_API_BASE}/me/services`;
  let method = 'POST';
  let body = {
    name: form.value.name,
    months: Number(form.value.selectedMonth),
    location: form.value.location,
    service: 'shadowsocks',
    prolong: Boolean(form.value.prolong),
    metadata: {
      plugin: form.value.selectedPlugin,
      method: form.value.selectedAlgorithm
    }
  };

  try {
    const response = await fetch(url, {
      method: method,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body)
    });

    if (!response.ok) {
      const errorData = await response.json();
      addToast({severity: 'error', summary: 'Ошибка', detail: errorData.error || 'Ошибка запроса', life: 3000});
    }

    switch (response.status) {
      case 401:
        console.error("Не авторизован");
        return;
      case 200:
        const text = await response.text();
        if (text) {
          const data = JSON.parse(text);
          window.location.href = `/#/services?id=${data.id}`
        }
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }
}

getLocationsFetch();

const isModalOpen = ref(false);
const openModal = () => {
  isModalOpen.value = true;
};

watch(isModalOpen, (isModalOpenNew) => {
  getMeFetch();
})


</script>

<template>
  <div>
    <div class="head">
      <h1 class="text-3xl font-extrabold">Создание SS</h1>
    </div>
    <form @submit.prevent="handleSubmit" class="flex flex-col gap-3 w-min mx-auto mt-3">
      <Modal :title="'Пополнение счёта'" :isOpen="isModalOpen" @update:isOpen="isModalOpen = $event">
        <div class="pl-4 border-l-4 border-gray-500 bg-gray-50 p-4">
          <p class="text-gray-700">
            Обычно, платёж зачисляется в течение минуты. Можете проверить в <a href="/#/account" target="_blank" class="font-semibold text-gray-600 underline hover:text-blue-800 transition duration-150">личном кабинете</a>.
          </p>
          <p class="mt-2 text-gray-600">
            Если деньги пришли, закройте это окно и нажмите «Оплатить».
          </p>
        </div>
        <div class="flex items-center gap-3 mt-4">
          <a href="/#/account?pay=true" target="_blank" class="bg-black text-white py-2 px-4 rounded-md shadow-lg transition duration-200">
            Перейти
          </a>
          <p class="text-sm text-gray-500">(откроется в новом окне)</p>
        </div>
      </Modal>

      <div class="input flex gap-2 items-center">
        <p class="w-24">Название:</p>
        <input type="text" v-model="form.name" placeholder="Название" class="p-2 border-2 rounded-md w-full"/>
      </div>
      <div class="input flex gap-2 items-center">
        <p class="w-24">Локация:</p>
        <div v-if="locations.length === 0">Загрузка локаций...</div>
        <select name="" id="" class="p-2 border-2 rounded-md w-full" v-model="form.location">
          <option value="none" disabled selected>Выберите локацию</option>
          <option v-for="location in locations" :key="location.name" :value="location.name">
            {{ location.name }}
          </option>
        </select>
      </div>
      <div class="input flex gap-2 items-center">
        <p class="w-24">Алгоритм:</p>
        <select v-model="form.selectedAlgorithm" class="p-2 border-2 rounded-md w-full">
          <option value="" disabled selected>Выберите алгоритм</option>
          <option v-for="algorithm in algorithms" :key="algorithm" :value="algorithm">
            {{ algorithm }}
          </option>
        </select>
      </div>
      <div class="input flex gap-2 items-center">
        <p class="w-24">Срок (мес):</p>
        <div class="flex space-x-2">
          <label class="flex items-center">
            <input
                type="radio"
                name="months"
                value="1"
                v-model="form.selectedMonth"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        1
      </span>
          </label>
          <label class="flex items-center">
            <input
                type="radio"
                name="months"
                value="3"
                v-model="form.selectedMonth"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        3
      </span>
          </label>
          <label class="flex items-center">
            <input
                type="radio"
                name="months"
                value="6"
                v-model="form.selectedMonth"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        6
      </span>
          </label>
          <label class="flex items-center">
            <input
                type="radio"
                name="months"
                value="12"
                v-model="form.selectedMonth"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        12
      </span>
          </label>
        </div>
      </div>
      <div class="input flex gap-2 items-center">
        <p class="w-24">Плагин:</p>
        <div class="flex space-x-2">
          <label class="flex items-center">
            <input
                type="radio"
                name="plugins"
                value="none"
                v-model="form.selectedPlugin"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        none
      </span>
          </label>
          <label class="flex items-center">
            <input
                type="radio"
                name="plugins"
                value="obfs"
                v-model="form.selectedPlugin"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        obfs
      </span>
          </label>
          <label class="flex items-center">
            <input
                type="radio"
                name="plugins"
                value="v2ray"
                v-model="form.selectedPlugin"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        v2ray
      </span>
          </label>
        </div>
      </div>
      <div class="input flex gap-2 items-center">
        <p class="w-24">Авт. продление:</p>
        <div class="flex space-x-2">
          <label class="flex items-center">
            <input
                type="radio"
                name="prolong"
                value="false"
                v-model="form.prolong"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        Выкл
      </span>
          </label>
          <label class="flex items-center">
            <input
                type="radio"
                name="prolong"
                value="true"
                v-model="form.prolong"
                class="hidden peer"
            />
            <span
                class="py-2 px-4 border rounded cursor-pointer transition-colors duration-200 peer-checked:bg-black peer-checked:text-white peer-checked:border-transparent"
            >
        Вкл
      </span>
          </label>
        </div>
      </div>
      <div class="flex items-center justify-between">
        <p class="text-2xl">Итого:</p>
        <p class="text-3xl">{{ total }} ₽</p>
      </div>
      <button type="submit" class="bg-black text-white p-3 rounded-md w-full">Оплатить</button>
    </form>
  </div>
</template>

<style scoped>

</style>
