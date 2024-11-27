<script setup>
import {inject, ref} from 'vue';

const addToast = inject('addToast');

const form = {login: undefined, password: undefined, invite: undefined}

const showPassword = ref(false)

async function handleSubmit() {
  try {
    await login(form.login, form.password);
  } catch (error) {
    addToast({severity: 'error', summary: 'Ошибка', detail: error, life: 3000});
  }
}

async function login(username, password) {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE}/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username,
        password: password,
      }),
      credentials: 'include', // Куки будут отправлены и сохранены автоматически
    });

    if (!response.ok) {
      const errorData = await response.json(); // Получаем данные ошибки, если есть
      throw (errorData.error || 'Ошибка авторизации');
    }

    switch (response.status) {
      case 401:
        // Обработка неавторизованного доступа
        console.error("Не авторизован");
        return;
      case 200:
        window.location.reload();
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
    <div class="head">
      <h1 class="text-3xl font-extrabold">Вход</h1>
    </div>
    <form @submit.prevent="handleSubmit" class="flex flex-col gap-3 w-min mx-auto">
      <div class="input flex gap-2 items-center">
        <p class="w-24">Логин:</p>
        <input type="text" v-model="form.login" placeholder="Логин" class="p-2 border-2 rounded-md"/>
      </div>
      <div class="input flex gap-2 items-center relative">
        <p class="w-24">Пароль:</p>
        <input :type="showPassword?'text':'password'" v-model="form.password" placeholder="Пароль"
               class="p-2 border-2 rounded-md"/>
        <div class="absolute right-2 cursor-pointer select-none" @click="showPassword = !showPassword">👁</div>
      </div>
      <button type="submit" class="bg-black text-white p-3 rounded-md w-full">Войти</button>
    </form>
  </div>
</template>

<style scoped>
/* Ваши стили */
</style>
