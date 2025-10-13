# VK Video Reviews

Приложение, предоставляющее API для управления рецензиями на видео. Оно позволяет пользователям добавлять, просматривать и управлять рецензиями на видео.

## Установка

1. **Генерация Swagger:**

   ```bash
   make swag
   ```

2. **Docker для запуска базы данных:**

   ```bash
   make docker-up
   ```

3. **Запуск миграции:**

   ```bash
   make migrate DATABASE_URL=postgres://username:password@localhost:5432/vkvideo_reviews_db?sslmode=disable
   ```

4. **Запуск приложения:**

   ```bash
   go run main.go
   ```

   Приложение поднимается в `http://localhost:8080`.

## Makefile команды

- `make swag` - сгенерировать Swagger документацию.
- `make docker-up` - запустить контейнер с базой данных PostgreSQL.
- `make docker-down` - остановить контейнер с базой данных.
- `make docker-logs` - просмотреть логи базы данных.
- `make migrate` - выполнить миграции базы данных.
