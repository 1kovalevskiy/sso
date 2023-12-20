package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/1kovalevskiy/sso/internal/entity"
	error_ "github.com/1kovalevskiy/sso/internal/error"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type (
	Auth interface {
		Login(ctx context.Context, email string, password string, appID int) (string, error)
		RegisterNewUser(ctx context.Context, email string, pass string) (int, error)
	}

	AuthRepo interface {
		InsertUser(ctx context.Context, email string, passHash []byte) (int, error)
		GetUser(ctx context.Context, email string) (entity.User, error)
		GetApp(ctx context.Context, id int) (entity.App, error)
	}
)

type AuthUseCase struct {
	log  *slog.Logger
	repo AuthRepo
}

func New(
	log *slog.Logger,
	repo AuthRepo,
) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
		log:  log,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.repo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			a.log.Warn("user not found", error_.Err(err))

			return "", fmt.Errorf("%s: %w", op, error_.ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", error_.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", error_.Err(err))

		return "", fmt.Errorf("%s: %w", op, error_.ErrInvalidCredentials)
	}

	app, err := a.repo.GetApp(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := NewToken(user, app)
	if err != nil {
		a.log.Error("failed to generate token", error_.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthUseCase) RegisterNewUser(ctx context.Context, email string, pass string) (int, error) {
	const op = "Auth.RegisterNewUser"

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

func NewToken(user entity.User, app entity.App) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	duration := time.Duration(app.TTLHours) * time.Hour
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
