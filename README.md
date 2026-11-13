# VK Video Reviews

Приложение, предоставляющее API для управления рецензиями на видео. Оно позволяет пользователям добавлять, просматривать и управлять рецензиями на видео.

## Установка

Установите Go 1.25 или новее и Docker с Compose v2. Из корня проекта выполните:

```bash
docker compose up -d --wait
go run ./cmd/app
```

При наличии Make обе команды можно выполнить через `make run`. При запуске
приложение применяет миграции PostgreSQL и загружает фильмы из VK Видео в фоне.
Откройте [Swagger UI](http://localhost:8080/swagger/index.html).

Для локального запуска используется пароль Redis из `config.yaml`. Чтобы задать
свой пароль, установите `APP_REDIS_PASS` в окружении перед запуском Compose и Go.

## Makefile команды

- `make docker-up` — запустить PostgreSQL и Redis и дождаться их готовности.
- `make docker-down` — остановить контейнеры.
- `make docker-logs` — просмотреть логи PostgreSQL и Redis.
- `make migrate` — выполнить миграции вручную через Goose CLI с заданным `DATABASE_URL`.
- `make api-gen` — обновить сгенерированный код после изменения OpenAPI.
