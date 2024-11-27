import {defineConfig} from 'vite'
import Pages from 'vite-plugin-pages'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        Pages({
            dirs: 'src/pages',
        }),
        vue()
    ],
})
