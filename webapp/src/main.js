import {createApp} from 'vue'
import './style.css'
import App from './App.vue'
import router from './plugins/router'
import toastPlugin from './plugins/toastPlugin';

createApp(App).use(router).use(toastPlugin).mount('#app')
