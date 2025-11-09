package sqlite

import (
	"database/sql"
	"fmt"
	"url-shortener2/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const origin = "storage.sqlite.New"

	base, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", origin, err)
	}

	if err = base.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", origin, err)
	}
	//Мы в таблице храним исходные URL и их укороченные версии - alias'ы
	stmt, err := base.Prepare(`CREATE TABLE IF NOT EXISTS urls(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL)`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", origin, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", origin, err)
	}
	stmt.Close()

	stmt2, err := base.Prepare(`CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias)`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", origin, err)
	}
	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", origin, err)
	}
	stmt2.Close()
	return &Storage{base}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const origin = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO urls(url, alias) VALUES(?, ?)`)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", origin, err)
	}

	resp, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %w", origin, storage.ErrUrlExists)
		}
		return -1, fmt.Errorf("%s: %w", origin, err)
	}

	insertedId, err := resp.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: failed to get last id: %w", origin, err)
	}

	stmt.Close()

	return insertedId, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const origin = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM urls WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", origin, err)
	}

	var url string
	row := stmt.QueryRow(alias)
	err = row.Scan(&url)
	stmt.Close()
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", origin, storage.ErrUrlNotFound)
		}
		return "", fmt.Errorf("%s: %w", origin, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const origin = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare(`DELETE FROM urls WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", origin, err)
	}

	res, _ := stmt.Exec(alias)
	stmt.Close()
	if num, _ := res.RowsAffected(); num == 0 {
		return storage.ErrUrlNotDeleted
	}

	return nil
}
