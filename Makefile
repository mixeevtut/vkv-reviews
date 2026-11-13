.PHONY: docker-up docker-down docker-logs migrate api-gen generate run

docker-up:
	docker compose up -d --wait postgres redis

docker-down:
	@echo "Stopping PostgreSQL database..."
	docker compose down

docker-logs:
	@echo "Database logs:"
	docker compose logs postgres redis

migrate:
	@echo "Running migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is required"; \
		echo "Example: make migrate DATABASE_URL=postgres://postgres:postgres@localhost:5432/vkvideo_reviews_db"; \
		exit 1; \
	fi
	goose -dir migrations postgres "$(DATABASE_URL)" up

OPENAPI=internal/api/openapi.yml
CFG=internal/api/config.yml

api-gen:
	go tool oapi-codegen -config $(CFG) $(OPENAPI)

generate: api-gen

run: docker-up
	@echo "Start app"
	go run ./cmd/app
