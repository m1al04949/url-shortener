package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/m1al04949/url-shortener/internal/storage"
	"modernc.org/sqlite"
	sqlerr "modernc.org/sqlite/lib"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s", storagePath))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// Migrator
	// stmt, err := db.Prepare(`
	// CREATE TABLE IF NOT EXISTS url_alias(
	// 	id INTEGER PRIMARY KEY,
	// 	alias TEXT NOT NULL UNIQUE,
	// 	url TEXT NOT NULL);
	// CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	// `)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// _, err = stmt.Exec()
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url_alias(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr *sqlite.Error

		if errors.As(err, &sqliteErr) {
			switch sqliteErr.Code() {
			case sqlerr.SQLITE_CONSTRAINT_UNIQUE:
				return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
			default:
				return 0, fmt.Errorf("%s: sqlite error [%d]: %w", op, sqliteErr.Code(), err)
			}
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url_alias WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement %w", op, err)
	}

	return resURL, nil

}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url_alias WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		var sqliteErr *sqlite.Error

		if errors.As(err, &sqliteErr) {
			switch sqliteErr.Code() {
			case sqlerr.SQLITE_CONSTRAINT_UNIQUE:
				return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
			default:
				return fmt.Errorf("%s: sqlite error [%d]: %w", op, sqliteErr.Code(), err)
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
