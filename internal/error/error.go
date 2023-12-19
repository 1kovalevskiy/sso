package error

import (
	"errors"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
