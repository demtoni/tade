<script setup>
const form = {login: undefined, password: undefined, invite: undefined}

function handleSubmit() {

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
      <div class="input flex gap-2 items-center">
        <p class="w-24">Пароль:</p>
        <input type="password" v-model="form.password" placeholder="Пароль" class="p-2 border-2 rounded-md"/>
      </div>
      <div class="input flex gap-2 items-center">
        <p class="w-24">Инвайт:</p>
        <input type="password" v-model="form.invite" placeholder="Инвайт" class="p-2 border-2 rounded-md"/>
      </div>
      <button type="submit" class="bg-black text-white p-3 rounded-md w-full">Войти</button>
    </form>
  </div>
</template>

<style scoped>

</style>