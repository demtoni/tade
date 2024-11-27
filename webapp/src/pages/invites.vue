<script setup>
import {ref} from "vue";

const invites = ref([])


async function newInvite() {
  let url = `${import.meta.env.VITE_API_BASE}/me/invites`;
  let method = 'POST';

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
        await getInvitesFetch()
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }
}

async function getInvitesFetch() {
  let url = `${import.meta.env.VITE_API_BASE}/me/invites`;
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
          invites.value = JSON.parse(text);
        }
        break;
      default:
        console.error("Ошибка ответа", response.status);
    }
  } catch (err) {
    console.error("Невозможно отправить запрос", err);
  }
}

getInvitesFetch()
const domainWithPort = window.location.host;
</script>

<template>
  <div>
    <div class="head flex items-center justify-between">
      <h1 class="text-3xl font-extrabold">Приглашения</h1>
      <button class="bg-black text-white p-2 rounded-md" @click="newInvite">➕ Создать</button>
    </div>
    <div class="content mt-3">
      <div v-for="invite in invites" class="p-3 bg-gray-100 rounded-md">
        <p>Ссылка для приглашения:</p>
        <textarea name="" id="" cols="30" rows="4"
                  class="resize-none bg-gray-900 p-2 rounded-md text-blue-200 cursor-text" disabled>{{domainWithPort}}/#/invite?code={{invite.invite_code}}</textarea>
      </div>
    </div>

  </div>
</template>

<style scoped>

</style>