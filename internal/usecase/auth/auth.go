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

type Token = string

func New(repository repository.Users) Auth {
	return Auth{
		repository: repository,
	}
}

var _ usecase.Auth = (*Auth)(nil)

func (a Auth) SignUp(ctx context.Context, username, password string) (Token, error) {
	passwordHash := password // TODO
	id, err := a.repository.CreateUser(ctx, username, passwordHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user")
	}

	token := *id // TODO

	return token.String(), nil
}

func (a Auth) SignIn(ctx context.Context, username, password string) (Token, error) {
	passwordHash := password // TODO
	id, err := a.repository.CheckUser(ctx, username, passwordHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to login")
	}

	token := *id // TODO

	return token.String(), nil
}
