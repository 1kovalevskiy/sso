package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/1kovalevskiy/sso/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Auth interface {
		GetCreateApp(ctx context.Context, name string, password string, secret string, ttlHour int) (int, error)
		Login(ctx context.Context, email string, password string, appID int) (string, error)
		RegisterNewUser(ctx context.Context, email string, pass string) (int, error)
	}

	AuthRepo interface {
		InsertUser(ctx context.Context, email string, passHash []byte) (int, error)
		GetUser(ctx context.Context, email string) (entity.User, error)
		GetAppForUser(ctx context.Context, id int) (entity.App, error)
		GetAppByName(ctx context.Context, name string) (entity.App, error)
		InsertApp(ctx context.Context, name string, passHash []byte, secret string, ttlHour int) (int, error)
		UpdateApp(ctx context.Context, id_ int, secret string, ttlHour int) (int, error)
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
