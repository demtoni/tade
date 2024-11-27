// src/hooks/useToast.js

import { reactive, ref } from 'vue';

// Реактивный список уведомлений
const toasts = ref([]);

// Функция для добавления уведомления
const addToast = ({ severity, summary, detail, life = 3000 }) => {
    const toast = reactive({
        summary,
        detail,
        severityClass: getSeverityClass(severity),
        visible: true,
    });

    // Ограничиваем количество тостов до 5
    if (toasts.value.length >= 5) {
        toasts.value.shift(); // Удаляем первый (самый старый) тост
    }

    // Добавляем новое уведомление в массив
    toasts.value.push(toast);

    // Убираем уведомление после указанного времени жизни (life)
    setTimeout(() => {
        toast.visible = false;
        // Удаляем уведомление после анимации исчезновения
        setTimeout(() => {
            toasts.value = toasts.value.filter(t => t !== toast);
        }, 300); // Задержка для завершения анимации
    }, life);
};

// Маппинг severity на классы Tailwind
const getSeverityClass = (severity) => {
    switch (severity) {
        case 'success':
            return 'bg-green-500';
        case 'info':
            return 'bg-blue-500';
        case 'warn':
            return 'bg-orange-500';
        case 'error':
            return 'bg-red-500';
        case 'secondary':
            return 'bg-gray-200 text-gray-800';
        case 'contrast':
            return 'bg-black text-white';
        default:
            return 'bg-green-500';
    }
};

// Экспортируем функцию useToast
export const useToast = () => {
    return { addToast, toasts };
};
