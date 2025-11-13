## обновление зависимостей
tidy:
	go mod tidy

migrate-up:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations