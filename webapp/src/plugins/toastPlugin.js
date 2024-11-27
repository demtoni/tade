import { useToast } from '../hooks/useToast';

export default {
    install(app) {
        const { addToast, toasts } = useToast();
        // Добавляем через provide
        app.provide('addToast', addToast);
        app.provide('toasts', toasts);  // Также предоставляем toasts, если нужно рендерить список тостов
    }
};