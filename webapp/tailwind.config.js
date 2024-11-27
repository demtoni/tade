/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/pages/**/*.{html,js,jsx,ts,tsx,vue}",
        "./src/components/**/*.{html,js,jsx,ts,tsx,vue}",
        "./src/**/*.{html,js,jsx,ts,tsx,vue}",
    ],
    theme: {
        extend: {
            screens: {
                '2xl': '1920px'
            },
        },
    },
    plugins: [],
}
