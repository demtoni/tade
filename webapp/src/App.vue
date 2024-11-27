<script setup>
import {RouterView} from "vue-router";
import Toast from "./components/Toast.vue";
import SideBar from "./components/sideBar.vue";
import {onMounted, onUnmounted, ref} from "vue";
import NavBar from "./components/navBar.vue";
import TopProfileInfo from "./components/topProfileInfo.vue";
import {useRoute} from 'vue-router';

const route = useRoute();
import {watchEffect} from 'vue';

const isNavOpen = ref(false)

const hideHeader = ref(true);

// Следим за изменением маршрута
watchEffect(() => {
  hideHeader.value = route.path === '/' || route.path === '/invite';
});

// Функция для обновления isNavOpen в зависимости от ширины экрана
const updateNavState = () => {
  isNavOpen.value = window.matchMedia("(max-width: 640px)").matches
}

onMounted(() => {
  updateNavState();
  // Установим isNavOpen в false при первой загрузке, если ширина экрана меньше 640px
  if (window.matchMedia("(max-width: 640px)").matches) {
    isNavOpen.value = false;
  }
  window.addEventListener('resize', updateNavState);
});


onUnmounted(() => {
  window.removeEventListener('resize', updateNavState)
})
</script>

<template>
  <div class="flex sm:flex-row mx-auto max-w-screen-2xl flex-col p-3">
    <Toast/>
    <sideBar :enabled="isNavOpen" title="Меню" position="left" @update:enabled="isNavOpen = $event">
      <navBar/>
    </sideBar>
    <header v-if="!hideHeader" class="sm:hidden flex items-center justify-between">
      <div class="bg-gray-200 w-fit p-2 rounded-full cursor-pointer" @click="isNavOpen=true">
        <img class="w-4 h-4" src="../src/assets/icons/menu.svg" alt="">
      </div>
      <div class="">
        <topProfileInfo/>
      </div>
    </header>
    <div class="hidden sm:block">
      <navBar/>
    </div>
    <router-view class="w-full mt-3"/>
  </div>
</template>

<style scoped>

</style>
