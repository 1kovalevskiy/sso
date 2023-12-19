package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_defaultTimeout = time.Second
)

type SQLite struct {
	Timeout time.Duration
	DB      *sql.DB
}

func New(url string, timeout string) (*SQLite, error) {
	db, err := openDB(url)
	if err != nil {
		return nil, err
	}
	to, err := time.ParseDuration(timeout)
	if err != nil {
		to = _defaultTimeout
	}

	mysql := &SQLite{
		Timeout: to,
		DB:      db,
	}

	return mysql, nil
}

func (p *SQLite) Close() {
	if p.DB != nil {
		p.DB.Close()
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
