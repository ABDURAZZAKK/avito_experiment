# Сервис динамического сегментирования пользователей

Выполнено основное и 1 Доп. Задание

Все эндпоинты описаны в экспортированном из Postman файле: `avito_exp.postman_collection.json`

Используемые технологии:
- PostgreSQL
- Docker
- Echo (фраймворк)
- golang-migrate/migrate
- pgx (драйвер PostgreSQL)
- RabbitMQ (для асинхронного создания CSV файлов)


## Запуск 

```bash
docker compose up 
```

