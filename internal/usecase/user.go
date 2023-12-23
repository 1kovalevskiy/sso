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

func (a *AuthUseCase) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "usecase.Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.repo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			log.Warn("user not found", error_.Err(err))

			return "", fmt.Errorf("%s: %w", op, error_.ErrInvalidCredentials)
		}

		log.Error("failed to get user", error_.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", error_.Err(err))

		return "", fmt.Errorf("%s: %w", op, error_.ErrInvalidCredentials)
	}

	app, err := a.repo.GetAppForUser(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := NewToken(user, app)
	if err != nil {
		log.Error("failed to generate token", error_.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthUseCase) RegisterNewUser(ctx context.Context, email string, pass string) (int, error) {
	const op = "usecase.Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", error_.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.repo.InsertUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", error_.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
