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

dev-start:
	docker compose -f ../compose-base.yml -f ../compose-dev.yml up -d

dev-stop:
	docker compose -f ../compose-base.yml -f ../compose-dev.yml down

dev-logs:
	docker compose -f ../compose-base.yml -f ../compose-dev.yml logs -f