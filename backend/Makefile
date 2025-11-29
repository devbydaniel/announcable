include ../.env
export

MIGRATION_DIR=internal/database/migrations
POSTGRES_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=disable

migrations-new:
	migrate create -ext sql -dir ${MIGRATION_DIR} -seq $(name)

migrations-up:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} up 1

migrations-up-all:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} up

migrations-down:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} down 1

migration-down-all:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} down

migration-force:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} force $(version)

migrations-unfuck:
	@echo "Getting current version..."
	@V=$$(migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} version); \
	echo "Forcing to version $$V..."; \
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} force $$V; \
	echo "Rolling back one migration..."; \
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} down 1

stripe-webhook:
	stripe listen --forward-to http://localhost:$(PORT)/stripe/webhook

# Development workflow:
# Terminal 1: make dev-services (start postgres, mail, minio, etc.)
# Terminal 2: make dev-air (Go hot-reload with env loaded)
# Terminal 3: npm run dev (Vite CSS/JS watch)

dev-services:
	docker compose -f ../compose-base.yml -f ../compose-dev.yml up -d --scale app=0

dev-air:
	@echo "Loading environment and starting Air..."
	@set -a; source ../.env; [ -f .env.local ] && source .env.local; set +a; air

dev-stop:
	docker compose -f ../compose-base.yml -f ../compose-dev.yml down

dev-logs:
	docker compose -f ../compose-base.yml -f ../compose-dev.yml logs -f

# Legacy command - kept for compatibility
dev-start: dev-services