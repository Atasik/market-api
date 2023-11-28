# Прототип API интернет-магазина
*Это немного видоизмененный бэкенд, поэтому CORS-ов тут нет, да и без фронта он не нужен

К сожалению, пока нет времени добавить вариант использования приложения без внешних сервисов (Хотя интерфейсы все прописаны). Поэтому дела обстоят так, что для работы добавления товаров, нужно обязательно зарегистрироваться в сервисе [Cloudinary](https://cloudinary.com/) (на 2023 год можно только через VPN) и занести следующие переменные окружения для работы приложения:
```````
CLOUDINARY_CLOUD=<ваш cloud из сервиса Cloudinary>
CLOUDINARY_KEY=<ваш key из сервиса Cloudinary>
CLOUDINARY_SECRET=<ваш secret из сервиса Cloudinary>

POSTGRES_PASSWORD=qwerty
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres

JWT_SIGNING_KEY=<ваш signing key>

HTTP_HOST=localhost
```````

Запуск:
```
make run
```

Чтобы протестировать api, надо зайти по этому адресу (если HTTP_HOST=localhost):
```
http://localhost:8080/swagger
```

Подробнее прочитать, где можно найти api-ключи, описанные выше, можно [здесь](https://cloudinary.com/documentation/admin_api#:~:text=Your%20Cloudinary%20API%20Key%20and,are%20used%20for%20the%20authentication.).

Road-map:
- [x] Регистрация, авторизация (пароли хэшируются);
- [x] Добавление, удаление, изменение, просмотр товара;
- [x] Возможность оставлять отзывы на продукты;
- [x] Добавление товара в корзину с дальнейшей возможностью покупки;
- [x] Хранение истории заказов;
- [x] Сортировка товаров по дате добавления/популярности;
- [ ] Возможность добавлять типы продуктов;
- [ ] unit-тесты (2%);
- [ ] JS-фронтенд;
- [x] Работающий Dockerfile;
- [ ] CI/CD;
- [ ] Рефакторинг некоторой части кода.
