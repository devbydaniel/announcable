include ../.env
export

MIGRATION_DIR=internal/database/migrations
POSTGRES_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_NAME)?sslmode=disable

migrations-new:
	migrate create -ext sql -dir ${MIGRATION_DIR} -seq $(name)

migrations-up:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} up

migrations-down:
	migrate -database ${POSTGRES_URL} -path ${MIGRATION_DIR} down
