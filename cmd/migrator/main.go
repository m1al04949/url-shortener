package main

import (
	"errors"
	"flag"
	"fmt"

	// Библиотека миграций
	"github.com/golang-migrate/migrate/v4"
	// Используем драйвер для modernc.org/sqlite
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	// Драйвер для файловых миграций
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// Регистрируем драйвер SQLite
	_ "modernc.org/sqlite"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
