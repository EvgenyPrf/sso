## обновление зависимостей
tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations

migrate-test:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test