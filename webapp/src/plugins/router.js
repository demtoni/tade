import routes from '~pages'
import {createRouter, createWebHashHistory} from "vue-router";

const router = createRouter({
    routes,
    history: createWebHashHistory()
});

// Глобальная навигационная защита
router.beforeEach((to, from, next) => {
    // Ищем куки с именем 'session'
    const sessionCookie = document.cookie.split('; ').find(row => row.startsWith('session='));

    // Разрешаем доступ к страницам invite/{string} только неавторизованным пользователям
    if (to.path.startsWith('/invite') && to.path !== '/invites') {
        if (sessionCookie) {
            next({path: '/home'}); // Если пользователь авторизован, перенаправляем на страницу home
        } else {
            next(); // Если не авторизован, продолжаем к invite
        }
    } else if (sessionCookie && to.path === '/') {
        next({path: '/home'}); // Перенаправляем на страницу home, если сессия существует
    } else if (!sessionCookie && to.path !== '/') {
        next({path: '/'}); // Если куки нет и пользователь пытается зайти на любую страницу, кроме индекса
    } else {
        next(); // Если ничего из вышеперечисленного не происходит, продолжаем
    }
});

export default router;
