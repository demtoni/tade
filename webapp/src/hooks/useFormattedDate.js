export function useFormattedDate() {
    function getDate(timestamp) {
        const date = new Date(timestamp * 1000);
        return date.toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            hour12: false // Использовать 24-часовой формат
        });
    }

    return { getDate };
}
