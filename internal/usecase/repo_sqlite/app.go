package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/1kovalevskiy/sso/internal/entity"
	"github.com/mattn/go-sqlite3"
)

func (r *AuthRepo) GetAppForUser(ctx context.Context, id int) (entity.App, error) {
	const op = "internal - usecase - repo_sqlite - AuthRepo.GetAppForUser"

	stmt, err := r.DB.Prepare(`SELECT id, name, pass_hash, secret, ttl_hours FROM apps WHERE id = ?`)
	if err != nil {
		return entity.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app entity.App
	err = row.Scan(&app.ID, &app.Name, &app.PassHash, &app.Secret, &app.TTLHours)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.App{}, fmt.Errorf("%s: %w", op, entity.ErrAppNotFound)
		}

		return entity.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil

}

func (r *AuthRepo) GetAppByName(ctx context.Context, name string) (entity.App, error) {
	const op = "internal - usecase - repo_sqlite - AuthRepo.GetAppByName"

	stmt, err := r.DB.Prepare(`SELECT id, name, pass_hash, secret, ttl_hours FROM apps WHERE name = ?`)
	if err != nil {
		return entity.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, name)

	var app entity.App
	err = row.Scan(&app.ID, &app.Name, &app.PassHash, &app.Secret, &app.TTLHours)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.App{}, fmt.Errorf("%s: %w", op, entity.ErrAppNotFound)
		}

		return entity.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil

}

func (r *AuthRepo) InsertApp(ctx context.Context, name string, passHash []byte, secret string, ttlHour int) (int, error) {
	const op = "internal - usecase - repo_sqlite - AuthRepo.InsertApp"

	stmt, err := r.DB.Prepare(`INSERT INTO apps(name, pass_hash, secret, ttl_hours) VALUES(?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, name, passHash, secret, ttlHour)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, entity.ErrAppExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

func (r *AuthRepo) UpdateApp(ctx context.Context, id_ int, secret string, ttlHour int) (int, error) {
	const op = "internal - usecase - repo_sqlite - AuthRepo.UpdateApp"

	stmt, err := r.DB.Prepare(`UPDATE apps SET secret = ?, ttl_hours = ? WHERE id = ?`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, secret, ttlHour, id_)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}

