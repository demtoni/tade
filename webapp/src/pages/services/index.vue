<script setup>
import {useFormattedDate} from '../../hooks/useFormattedDate';
import {useRoute} from 'vue-router';
import {ref} from "vue";
import QRCode from 'qrcode';
import Modal from "../../components/modal.vue";

const isModalOpen = ref(false);

const url = ref(null);
const route = useRoute();
const serviceId = route.query.id;
const qrCodeUrl = ref('');

const {getDate} = useFormattedDate();

const serviceData = ref({})

// Функция для открытия модального окна
const openModal = () => {
  isModalOpen.value = true;
  generateQRCode();
};

const generateQRCode = async () => {
  qrCodeUrl.value = await qrCodeGenerate(serviceData.value.metadata.connect_url);
};

// Функция для генерации QR-кода
async function qrCodeGenerate(text) {
  try {
    return await QRCode.toDataURL(text);
  } catch (err) {
    console.error(err);
    return null; // Возвращаем null в случае ошибки
  }
}

async function getServiceData() {
  let url = `${import.meta.env.VITE_API_BASE}/me/services/${serviceId}`;
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
          serviceData.value = JSON.parse(text);
        }
        break;
      default:
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }
}

getServiceData();
</script>

<template>
  <div>
    <div class="head">
      <h1 class="text-3xl font-extrabold">{{ serviceData.name }}</h1>
    </div>
    <div class="content mt-3">
      <div class="flex p-2 items-center gap-3" :class="serviceData.name?'':'skeleton'">
        <div class="bg-gray-100 rounded-lg m-1 flex flex-col items-start gap-4 p-3 overflow-hidden">
          <div class="flex w-full gap-12 flex-wrap">
            <div class="flex flex-col items-start">
              <p class="text-gray-500">Истекает после</p>
              <p>{{ getDate(serviceData.expires_at) }}</p>
            </div>
            <div class="flex flex-col items-start">
              <p class="text-gray-500">Сервис</p>
              <p>{{ serviceData.service }}</p>
            </div>
            <div class="flex flex-col items-start">
              <p class="text-gray-500">Локация</p>
              <p>{{ serviceData.location }}</p>
            </div>
            <div class="flex flex-col items-start" v-if="serviceData && serviceData.metadata">
              <div class="text-gray-500 mb-3 flex items-center">
                <p>Ссылка для подключения:</p>
              </div>
              <textarea name="" id="" cols="30" rows="4"
                        class="resize-none bg-gray-900 p-2 rounded-md text-blue-200 cursor-text" disabled>{{ serviceData.metadata.connect_url }}</textarea>
            </div>
            <div class="flex flex-col items-start" v-if="serviceData && serviceData.metadata">
              <div class="text-gray-500 mb-3 flex items-center">
                <p>QR-код для подключения:</p>
              </div>
              <button class="bg-black text-white p-2 rounded-md" @click="openModal">Показать</button>
              <Modal :title="'QR-Code'" :isOpen="isModalOpen" @update:isOpen="isModalOpen = $event">
                <img class="mx-auto" v-if="qrCodeUrl" :src="qrCodeUrl" alt="QR Code">
              </Modal>
            </div>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
