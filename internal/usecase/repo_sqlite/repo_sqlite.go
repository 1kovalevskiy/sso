package repo

import (
	"github.com/1kovalevskiy/sso/pkg/sqlite"
)

type AuthRepo struct {
	*sqlite.SQLite
}

func New(mysql_ *sqlite.SQLite) *AuthRepo {
	return &AuthRepo{mysql_}
}
