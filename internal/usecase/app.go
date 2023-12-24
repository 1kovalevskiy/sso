package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/1kovalevskiy/sso/internal/entity"
	error_ "github.com/1kovalevskiy/sso/internal/error"
	"golang.org/x/crypto/bcrypt"
)


func (a *AuthUseCase) createApp(ctx context.Context, name string, pass string, secret string, ttlHour int) (int, error) {
	const op = "internal - usecase - Auth.createApp"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", name),
	)

	log.Info("registering app")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", error_.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.repo.InsertApp(ctx, name, passHash, secret, ttlHour)
	if err != nil {
		log.Error("failed to save app", error_.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *AuthUseCase) getAppByName(ctx context.Context, name string, password string) (int, error) {
	const op = "internal - usecase - Auth.getAppByName"

	log := a.log.With(
		slog.String("op", op),
		slog.String("service_name", name),
	)

	log.Info("attempting to get app")

	app, err := a.repo.GetAppByName(ctx, name)
	if err != nil {
		if errors.Is(err, entity.ErrAppNotFound) {
			log.Warn("app not found", error_.Err(err))
			return 0, entity.ErrAppNotFound
		}

		log.Error("failed to get app", error_.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(app.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", error_.Err(err))

		return 0, fmt.Errorf("%s: %w", op, error_.ErrInvalidCredentials)
	}
	log.Info("app successfully identify")


	return app.ID, nil

}

func (a *AuthUseCase) updateApp(ctx context.Context, id_ int, secret string, ttlHour int) (int, error) {
	const op = "internal - usecase - Auth.updateApp"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("service_id", id_),
	)

	log.Info("attempting to update app")

	id, err := a.repo.UpdateApp(ctx, id_, secret, ttlHour)
	if err != nil {
		log.Error("failed to save app", error_.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (a *AuthUseCase) GetCreateApp(ctx context.Context, name string, password string, secret string, ttlHour int) (int, error) {

	id, err := a.getAppByName(ctx, name, password)
	if err != nil && errors.Is(err, entity.ErrAppNotFound) {
		return a.createApp(ctx, name, password, secret, ttlHour)
	}
	if err != nil {
		return 0, err
	}

	return a.updateApp(ctx, id, secret, ttlHour)

}
