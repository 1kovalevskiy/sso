package entity

import "errors"

var (
	ErrAppExists   = errors.New("app already exists")
	ErrAppNotFound = errors.New("app not found")
)

type App struct {
	ID			int
	Name		string
	PassHash	[]byte
	Secret		string
	TTLHours	int
}
