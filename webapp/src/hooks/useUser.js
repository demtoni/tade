import {ref} from 'vue';
import Cookies from 'js-cookie';

export function useUser() {
    const meData = ref(null);

    // Функция для получения данных о пользователе
    const getMeFetch = async () => {
        let url = `${import.meta.env.VITE_API_BASE}/me`;
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
                        meData.value = JSON.parse(text);
                    }
                    break;
                default:
                    console.error("Ошибка ответа", response.status);
            }
        } catch (err) {
            console.error("Невозможно отправить запрос", err);
        }
    };

    // Функция для выхода (logout)
    const logout = async () => {
    	Cookies.remove('session', { path: '/', domain: `.${import.meta.env.VITE_DOMAIN}` });
    	window.location.reload();
    };

    return {meData, getMeFetch, logout};
}
