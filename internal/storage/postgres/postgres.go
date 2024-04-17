package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connectionsString string) (*Storage, error) {
	const op = "storage.postgres.New" // Путь до функции
	const firstPreparedQuery = `
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);`
	const secondPreparedQuery = `CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`

	db, err := sql.Open("postgres", connectionsString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = preparePostgres(db, firstPreparedQuery, op); err != nil {
		return nil, err
	}

	if err = preparePostgres(db, secondPreparedQuery, op); err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func preparePostgres(db *sql.DB, query string, op string) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"
	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resURL, nil
}
