.PHONY: local local-log dev dev-log down migrate-create migrate-up migrate-down db-connect

local:
	docker compose -f docker-compose.local.yml down
	docker compose -f docker-compose.local.yml rm -f
	docker compose -f docker-compose.local.yml build
	docker compose -f docker-compose.local.yml up -d
	docker compose -f docker-compose.local.yml ps

local-log:
	docker compose -f docker-compose.local.yml logs -f app db

down:
	docker compose -f docker-compose.local.yml down
	docker compose -f docker-compose.local.yml rm -f

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>"; \
		exit 1; \
	fi
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	echo "-- +migrate Up\n\n-- +migrate Down" > db/migrations/$${timestamp}_$(name).sql; \
	echo "Created db/migrations/$${timestamp}_$(name).sql"

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

db-connect:
	docker compose -f docker-compose.local.yml exec db mysql -u app -papp app
