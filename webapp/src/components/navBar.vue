<script setup>
const token = document.cookie.split('; ').find(row => row.startsWith('session='));
</script>

<template>
  <div>
    <nav class="flex flex-col p-3 gap-3 w-64">
      <router-link v-if="token" to="/home" class="hover:underline text-xl ">Мои услуги</router-link>
      <router-link v-if="token" to="/create" class="hover:underline text-xl">Заказать услугу</router-link>
      <router-link v-if="token" to="/account" class="hover:underline text-xl">Личный кабинет</router-link>
      <router-link v-if="token" to="/invites" class="hover:underline text-xl">Приглашения</router-link>
      <router-link v-if="token" to="/support" class="hover:underline text-xl">Помощь</router-link>
      <router-link v-else to="/" class="hover:underline text-xl">Вход</router-link>
    </nav>
  </div>
</template>

<style scoped>

</style>