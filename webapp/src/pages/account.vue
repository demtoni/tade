<script setup>

import {ref} from "vue";
import {useFormattedDate} from '../hooks/useFormattedDate';
import {inject} from 'vue';
import SideBar from "../components/sideBar.vue";
import TopProfileInfo from "../components/topProfileInfo.vue";
import {useUser} from '../hooks/useUser.js';
import {onMounted} from 'vue';
import {useRoute} from "vue-router";

const route = useRoute();
const {meData, getMeFetch} = useUser();
const openPayPanel = Boolean(route.query.pay);

const meDataFetched = ref(false)

onMounted(() => {

  getMeFetch().then(() => {
    meDataFetched.value = true; // Устанавливаем в true после успешной загрузки данных
  });

  if (openPayPanel) {
    isPayOpen.value = true;
  }
});

const {getDate} = useFormattedDate();

const addToast = inject('addToast');
const transactionsData = ref({})
const transactionsDataLoaded = ref(false)
const passwordForm = ref({password: '', newPassword: '', confirmPassword: ''})

function setAmount(value) {
  amount.value = value;
}

async function handleSubmit() {
  try {
    await newPassword(passwordForm.value.password, passwordForm.value.newPassword, passwordForm.value.confirmPassword);
  } catch (error) {
    addToast({severity: 'error', summary: 'Ошибка', detail: error, life: 3000});
  }
}

async function newPassword(old_password, new_password, confirm_password) {
  if (new_password !== confirm_password) {
    addToast({severity: 'warn', summary: 'Внимание', detail: 'Пароли не совпадают', life: 3000});
    return
  }
  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE}/me/password`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        old_password: old_password,
        new_password: new_password,
      }),
      credentials: 'include', // Куки будут отправлены и сохранены автоматически
    });

    if (!response.ok) {
      const errorData = await response.json(); // Получаем данные ошибки, если есть
      throw (errorData.error || 'Ошибка изменения пароля');
    }

    switch (response.status) {
      case 401:
        console.error("Не авторизован");
        return;
      case 200:
        passwordForm.password = '';
        passwordForm.confirmPassword = '';
        addToast({severity: 'success', summary: 'Успех', detail: 'Пароль обновлен!', life: 3000});
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }

  } catch (error) {
    console.error('Ошибка:', error.message || error);
    throw error;
  }
}

async function getTransactions() {
  let url = `${import.meta.env.VITE_API_BASE}/me/transactions`;
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
          transactionsData.value = JSON.parse(text).reverse();
        }
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }
  transactionsDataLoaded.value = true;
}

const statusLabels = {
  completed: 'Выполнено',
  in_process: 'В процессе',
  canceled: 'Отменено',
};

getTransactions()

const isPayOpen = ref(false)
const amount = ref(100);
const selectedPayment = ref('bank_card');

async function handlePayment() {
  if (amount.value <= 0) {
    addToast({severity: 'error', summary: 'Ошибка', detail: "Сумма должна быть положительной", life: 3000});
  }

  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE}/me/balance`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        "amount": amount.value,
        "payment_method": selectedPayment.value,
        "return_url": import.meta.env.VITE_RETURN_URL
      }),
      credentials: 'include',
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw (errorData.error || 'Ошибка изменения пароля');
    }

    switch (response.status) {
      case 401:
        console.error("Не авторизован");
        return;
      case 200:
        const text = await response.text();
        if (text) {
          const data = JSON.parse(text);
          window.location.href = data.payment_url;
        }
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }

  } catch (error) {
    console.error('Ошибка:', error.message || error);
    throw error;
  }
}

</script>

<template>
  <div>
    <sideBar :enabled="isPayOpen" title="Пополнение" position="right" @update:enabled="isPayOpen = $event">
      <div class="flex flex-col items-center h-screen p-3 overflow-hidden">
        <form @submit.prevent="handlePayment" class="flex flex-col w-full gap-3">
          <div>
            <h2 class="mb-2 text-lg">Способ оплаты:</h2>
            <div class="space-y-2">
              <label
                  class="flex items-center p-2 border rounded-lg cursor-pointer transition-all"
                  :class="{
          'bg-gray-100 border-black': selectedPayment === 'bank_card',
          'border-gray-200 hover:bg-gray-50': selectedPayment !== 'bank_card'
        }"
              >
                <input
                    type="radio"
                    v-model="selectedPayment"
                    value="bank_card"
                    class="hidden"
                />
                <span class="flex items-center">
          <i class="fas fa-credit-card mr-2"></i> <!-- Иконка карты -->
          Банковские карты
        </span>
              </label>

              <label
                  class="flex items-center p-2 border rounded-lg cursor-pointer transition-all"
                  :class="{
          'bg-gray-100 border-black': selectedPayment === 'sbp',
          'border-gray-200 hover:bg-gray-50': selectedPayment !== 'sbp'
        }"
              >
                <input
                    type="radio"
                    v-model="selectedPayment"
                    value="sbp"
                    class="hidden"
                />
                <span class="flex items-center">
          <i class="fas fa-mobile-alt mr-2"></i> <!-- Иконка для СБП -->
          СБП
        </span>
              </label>
            </div>
          </div>

          <div>
            <label for="amount" class="mb-2 text-lg">Сумма оплаты:</label>
            <input
                type="number"
                id="amount"
                v-model="amount"
                class="border border-gray-300 rounded p-2 mb-4"
                placeholder="Введите сумму"
                min="1"
                required
            />
            <div class="flex gap-2 mt-2 flex-wrap">
              <button
                  @click.prevent="setAmount(100)"
                  class="px-4 py-2 bg-gray-100 border rounded-lg hover:bg-gray-200"
              >
                100
              </button>
              <button
                  @click.prevent="setAmount(300)"
                  class="px-4 py-2 bg-gray-100 border rounded-lg hover:bg-gray-200"
              >
                300
              </button>
              <button
                  @click.prevent="setAmount(600)"
                  class="px-4 py-2 bg-gray-100 border rounded-lg hover:bg-gray-200"
              >
                600
              </button>
            </div>
          </div>
          <button
              type="submit"
              class="bg-black text-white rounded p-2 hover:bg-gray-900"
          >
            Оплатить
          </button>
        </form>
      </div>
    </sideBar>
    <div class="head flex items-center justify-between">
      <h1 class="text-3xl font-extrabold">Личный кабинет</h1>
      <div class="hidden sm:block">
        <topProfileInfo/>
      </div>
    </div>
    <div class="content mt-3">
      <div class="block items-stretch gap-3 lg:flex">
        <div class="bg-gray-100 rounded-lg m-1 flex flex-col items-start gap-4 p-3 w-full"
             :class="meDataFetched?'':'skeleton'">
          <h2 class="mb-2">Баланс</h2>
          <p class="md:text-6xl text-base font-semibold"><span v-if="meData">{{ meData.balance }} ₽</span></p>
          <button class="bg-black text-white p-2 rounded-md h-min" @click="isPayOpen = true">Пополнить баланс</button>
        </div>
        <div class="bg-gray-100 rounded-lg m-1 flex flex-col items-start gap-4 p-3 w-full">
          <h2 class="mb-2">Изменить пароль</h2>
          <form @submit.prevent="handleSubmit" class="flex flex-col gap-3">
            <div class="input flex flex-wrap gap-2 items-center">
              <p class="w-48">Текущий пароль:</p>
              <input type="password" v-model="passwordForm.password" name="currentPassword" placeholder="Текущий пароль"
                     class="p-2 border-2 rounded-md sm:w-auto w-full" minlength="8" maxlength="72" required
                     autocomplete="current-password"/>
            </div>
            <div class="input flex flex-wrap gap-2 items-center">
              <p class="w-48">Новый пароль:</p>
              <input type="password" v-model="passwordForm.newPassword" name="newPassword" placeholder="Новый пароль"
                     class="p-2 border-2 rounded-md sm:w-auto w-full" required autocomplete="new-password"/>
            </div>
            <div class="input flex flex-wrap gap-2 items-center">
              <p class="w-48">Подтверждение пароля:</p>
              <input type="password" v-model="passwordForm.confirmPassword" name="confirmPassword"
                     placeholder="Подтвердите пароль"
                     class="p-2 border-2 rounded-md sm:w-auto w-full" required/>
            </div>
            <button type="submit" class="bg-black text-white p-3 rounded-md w-full">Изменить</button>
          </form>
        </div>
      </div>
      <div>
        <h2 class="text-3xl font-bold">Transactions</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mt-3">
          <div v-for="item in [1,2,3]" :class="!transactionsData?'':'skeleton'" v-if="!transactionsDataLoaded">
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Aspernatur est eum ipsum placeat porro sapiente
            soluta temporibus. Beatae dolores ducimus fuga sequi voluptas voluptatibus! Aut culpa doloremque eius in
            ullam?
          </div>
          <div
              class="bg-white border border-gray-200 rounded-lg shadow-lg p-4 hover:shadow-xl transition-shadow duration-300 ease-in-out"
              v-for="transaction in transactionsData"
              :key="transaction.id"
          >
            <div class="flex justify-between">
              <h3 class="text-sm font-semibold text-gray-800">{{ getDate(transaction.timestamp) }}</h3>
              <p class="text-sm font-semibold text-gray-800">{{ transaction.amount }}₽</p>
            </div>
            <p class="mt-2 text-sm text-gray-600">
              Статус: <span class="font-bold">{{ statusLabels[transaction.status] || 'Неизвестный статус' }}</span>
            </p>
            <a
                class="mt-3 inline-block text-blue-600 hover:underline"
                :href="transaction.url"
            >
              Перейти
            </a>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
