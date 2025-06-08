package repository

import (
	"context"
	"errors"
)

var ErrorUserAlreadyExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrInvalidRequestFormat = errors.New("invalid request")
var ErrUnknownInternalError = errors.New("unknown internal error")

type Users interface {
	CreateUser(ctx context.Context, username, passwordHash string) error
	CheckUser(ctx context.Context, username, passwordHash string) error
}
