# Прототип интернет-магазина

Запуск
```
docker-compose up
go run cmd/market/main.go
```

Реализованный функционал:
- Регистрация, авторизация;
- Добавление, удаление, изменение, просмотр товара;
- Добавление товара в корзину с дальнейшей возможностью покупки;
- Хранение истории заказов;
- Сортировка товаров по дате добавления/популярности;

Road-map (Что ещё нужно реализовать):
- Хэширование паролей (!!!);
- Возможность пользователям оставлять отзывы на продукты;
- Возможность добавлять типы продуктов;
- Небольшой рефакторинг отдельных частей кода;
- Возможность смены языка;
- Доработать Dockerfile и, возможно, добавить CI/CD;
- Доработать дизайн;
- Тёмная/светлая тема;