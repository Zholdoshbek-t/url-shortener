package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	smtm, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_url_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer smtm.Close()
	_, err = smtm.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"
	smtm, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer smtm.Close()
	res, err := smtm.Exec(urlToSave, alias)
	if err != nil {
		if sqliteError, ok := err.(sqlite3.Error); ok && sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, errors.New("could not extract [id] from result"))
	}
	return id, nil
}

type URL struct {
	ID    int64
	Url   string
	Alias string
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	row := s.db.QueryRow("SELECT id, url, alias FROM url WHERE alias = ?", alias)
	var url URL
	err := row.Scan(&url.ID, &url.Url, &url.Alias)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url.Url, nil
}

func (s *Storage) GetAllURLs() ([]URL, error) {
	const op = "storage.sqlite.GetAllURLs"
	rows, err := s.db.Query("SELECT id, url, alias FROM url")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	urls := []URL{}
	for rows.Next() {
		var url URL
		rows.Scan(&url.ID, &url.Url, &url.Alias)
		urls = append(urls, url)
	}
	return urls, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"
	smtm, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer smtm.Close()
	_, err = smtm.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}
