# VK Video Reviews

Приложение, предоставляющее API для управления рецензиями на видео. Оно позволяет пользователям добавлять, просматривать и управлять рецензиями на видео.

## Установка

1. **Docker для запуска базы данных:**

   ```bash
   make docker-up
   ```

2. **Запуск приложения:**

   ```bash
   make run
   ```

3. **Запуск api-gen**
   
   ``` bash
   make api-gen 
   ```

   Приложение поднимается в `http://localhost:8080`.

   [SwaggerUI (localhost)](http://localhost:8080/swagger/index.html#/default)

## Makefile команды

- `make docker-up` - запустить контейнер с базой данных PostgreSQL.
- `make docker-down` - остановить контейнер с базой данных.
- `make docker-logs` - просмотреть логи базы данных.
- `make migrate` - выполнить миграции базы данных.
