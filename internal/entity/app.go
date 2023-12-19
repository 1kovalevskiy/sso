package entity

import "errors"

var ErrAppNotFound = errors.New("app not found")

type App struct {
	ID        int
	Name      string
	Secret    string
	TTL_hours int
}
