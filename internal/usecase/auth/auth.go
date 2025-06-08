package auth

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/pkg/errors"
)

type Auth struct {
	repository repository.Users
}

func New(repository repository.Users) Auth {
	return Auth{
		repository: repository,
	}
}

var _ usecase.Auth = (*Auth)(nil)

func (a Auth) SignUp(ctx context.Context, username, password string) (string, error) {
	passwordHash := password // TODO
	err := a.repository.CreateUser(ctx, username, passwordHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user")
	}

	token := "secret_token" // TODO

	return token, nil
}

func (a Auth) SignIn(ctx context.Context, username, password string) (string, error) {
	passwordHash := password // TODO
	err := a.repository.CheckUser(ctx, username, passwordHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to login")
	}

	token := "secret_token" // TODO

	return token, nil
}
