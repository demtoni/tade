<script setup>
import {inject, ref} from 'vue';
import {useRoute} from 'vue-router';

const route = useRoute();
const addToast = inject('addToast');
const invite = route.query.code;

const form = {login: undefined, password: undefined, invite: undefined}
const showPassword = ref(false)

async function handleSubmit() {
  try {
    await login(form.login, form.password, invite);
  } catch (error) {
    addToast({severity: 'error', summary: 'Ошибка', detail: error, life: 3000});
  }
}

async function login(username, password, invite) {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_BASE}/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username,
        password: password,
        invite: invite,
      }),
      credentials: 'include', // Куки будут отправлены и сохранены автоматически
    });

    if (!response.ok) {
      const errorData = await response.json(); // Получаем данные ошибки, если есть
      throw (errorData.error || 'Ошибка авторизации');
    }

    switch (response.status) {
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
      <h1 class="text-3xl font-extrabold">Регистрация</h1>
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
      <button type="submit" class="bg-black text-white p-3 rounded-md w-full">Зарегестрироваться</button>
    </form>
  </div>
</template>

<style scoped>
/* Ваши стили */
</style>
