.PHONY: swag docker-up docker-down docker-logs
swag:
	swag init -g cmd/app/main.go --parseDependency --parseInternal -o ./docs

docker-up:
	@echo "Starting PostgreSQL database..."
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@until docker-compose exec -T postgres pg_isready -U postgres -d vkvideo_reviews_db; do sleep 2; done
	@echo "Database is ready!"

docker-down:
	@echo "Stopping PostgreSQL database..."
	docker-compose down

docker-logs:
	@echo "Database logs:"
	docker-compose logs postgres

migrate:
	@echo "Running migrations..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Error: DATABASE_URL environment variable is required"; \
		echo "Example: make migrate DATABASE_URL=postgres://postgres:postgres@localhost:5432/vkvideo_reviews_db"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

OPENAPI=internal/api/openapi.yml
CFG=internal/api/config.yml

.PHONY: api-gen
api-gen:
	go tool oapi-codegen -config $(CFG) $(OPENAPI)

generate: api-gen