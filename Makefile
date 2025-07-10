STORAGE_PATH := ./storage/storage.db
MIGRATIONS_PATH := ./migrations

migrate:
	go run ./cmd/migrator --storage-path=$(STORAGE_PATH) \
		--migrations-path=$(MIGRATIONS_PATH)

start:
	go run ./cmd/url-shortener/main.go