const tg = window.Telegram.WebApp;

// Инициализация WebApp
tg.ready();
tg.expand();

// Устанавливаем тему
document.body.style.backgroundColor = tg.themeParams.bg_color || '#ffffff';

const cityInput = document.getElementById('cityInput');
const weatherBtn = document.getElementById('weatherBtn');
const status = document.getElementById('status');

// Обработка нажатия Enter в поле ввода
cityInput.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        getWeather();
    }
});

// Функция для отображения статуса
function showStatus(message, type) {
    status.innerHTML = message;
    status.className = 'status ' + type;
    status.style.display = 'block';
    
    // Скрываем статус через 5 секунд для успешных сообщений
    if (type === 'success') {
        setTimeout(() => {
            hideStatus();
        }, 5000);
    }
}

// Функция для скрытия статуса
function hideStatus() {
    status.style.display = 'none';
    status.className = 'status';
}

// Основная функция получения погоды
function getWeather() {
    const city = cityInput.value.trim();
    
    if (!city) {
        showStatus('❌ Пожалуйста, введите название города', 'error');
        cityInput.focus();
        return;
    }
    
    // Валидация ввода
    if (city.length < 2) {
        showStatus('❌ Название города слишком короткое', 'error');
        cityInput.focus();
        return;
    }
    
    if (city.length > 50) {
        showStatus('❌ Название города слишком длинное', 'error');
        cityInput.focus();
        return;
    }
    
    // Показываем индикатор загрузки
    weatherBtn.disabled = true;
    weatherBtn.innerHTML = '<span class="loading"></span>Получаем данные...';
    hideStatus();
    
    // Отправляем данные боту
    try {
        const data = {
            city: city,
            timestamp: Date.now(),
            user_id: tg.initDataUnsafe?.user?.id || null
        };
        
        tg.sendData(JSON.stringify(data));
        showStatus('✅ Запрос отправлен! Результат придет в чат.', 'success');
        
        // Очищаем поле ввода после успешной отправки
        cityInput.value = '';
        
        // Закрываем WebApp через 2 секунды
        setTimeout(() => {
            tg.close();
        }, 2000);
        
    } catch (error) {
        showStatus('❌ Ошибка отправки запроса. Попробуйте еще раз.', 'error');
        console.error('Error sending data:', error);
    }
    
    // Возвращаем кнопку в исходное состояние
    setTimeout(() => {
        weatherBtn.disabled = false;
        weatherBtn.innerHTML = '🌡️ Получить погоду';
    }, 2000);
}

// Дополнительные настройки WebApp
tg.MainButton.hide();
tg.BackButton.hide();

// Обработка событий WebApp
tg.onEvent('mainButtonClicked', function() {
    getWeather();
});

tg.onEvent('backButtonClicked', function() {
    tg.close();
});

// Настройка заголовка WebApp
if (tg.platform) {
    tg.setHeaderColor('#0088cc');
}

// Автофокус на поле ввода (для десктопа)
if (!tg.platform.includes('mobile')) {
    setTimeout(() => {
        cityInput.focus();
    }, 300);
}

// Логирование для отладки
console.log('WebApp initialized:', {
    platform: tg.platform,
    version: tg.version,
    user: tg.initDataUnsafe?.user,
    colorScheme: tg.colorScheme
});

// Обработка изменения темы
tg.onEvent('themeChanged', function() {
    document.body.style.backgroundColor = tg.themeParams.bg_color || '#ffffff';
});

// Функция для отображения информации о пользователе (для отладки)
function showUserInfo() {
    if (tg.initDataUnsafe?.user) {
        const user = tg.initDataUnsafe.user;
        console.log('User info:', user);
    }
}

// Запускаем функцию при загрузке
showUserInfo(); 