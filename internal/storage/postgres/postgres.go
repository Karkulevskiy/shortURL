package postgres

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connectionsString string) (*Storage, error) {
	const op = "storage.postgres.New" // Путь до функции
	const firstPreparedQuery = `
	CREATE TABLE IF NOT EXISTS url(
		id SERIAL PRIMARY KEY,
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

	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := stmt.Query(urlToSave, alias)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Constraint != "" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}
	if rows.Err() != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = $1`)
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

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"
	const query = `DELETE FROM url WHERE alias = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetAllURL() ([]string, error) {
	const op = "storage.postgres.GetAllURL"
	const query = `SELECT (url, alias) FROM url`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var url string
	var urls []string
	for rows.Next() {
		err = rows.Scan(&url)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		urls = append(urls, url)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return urls, nil
}

func (s *Storage) DropTable() error {
	const op = "storage.postgres.DropTable"
	const query = `DROP TABLE url`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
