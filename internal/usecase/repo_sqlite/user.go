package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/1kovalevskiy/sso/internal/entity"
	"github.com/mattn/go-sqlite3"
)

func (r *AuthRepo) InsertUser(ctx context.Context, email string, passHash []byte) (int, error) {
	const op = "internal.usecase.repo_sqlite.SaveUser"

	stmt, err := r.DB.Prepare(`INSERT INTO users(email, pass_hash) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, entity.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

func (r *AuthRepo) GetUser(ctx context.Context, email string) (entity.User, error) {
	const op = "internal.usecase.repo_sqlite.GetUser"

	stmt, err := r.DB.Prepare(`SELECT id, email, pass_hash FROM users WHERE email = ?`)
	if err != nil {
		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user entity.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("%s: %w", op, entity.ErrUserNotFound)
		}

		return entity.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil

}
