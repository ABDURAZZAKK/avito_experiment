# Сервис динамического сегментирования пользователей

Выполнено основное и Все Доп. Задания

Все эндпоинты описаны в экспортированном из Postman файле: `avito_exp.postman_collection.json`

Используемые технологии:
- PostgreSQL
- Docker
- Echo (фраймворк)
- golang-migrate/migrate
- pgx (драйвер PostgreSQL)
- RabbitMQ (для асинхронного создания CSV файлов и отложенных задач)


## Запуск 

```bash
docker compose up 
```

Возникшие в ходе выполнения вопросы и ответы на них:

>1 Доп задание сохранение статистики попадания или удалиниия пользователя из сегмента.
>В статистике должно быть отраженно удаление пользователя или сегмента из БД ?
>
>Нет.

>2 Доп задание "реализовать возможность задавать TTL, для этого в метод добавления сегментов пользователю передаём время удаления пользователя из сегмента отдельным полем".
>Сегменты из какого списка (на добавление или удаление) должны попасть под отложенное удаление ?
>
>Будут удалятся сегменты из списка на добавление.

>3 Доп задание. Нужно ли при создании сегментов с автоматическим добавлением пользователю, добавить возможность отложенного удаления сегментов у пользователей ?
>
>Да.