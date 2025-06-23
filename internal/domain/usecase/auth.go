package usecase

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
)

type Auth interface {
	SignUp(ctx context.Context, username, password string) (string, error)
	SignIn(ctx context.Context, username, password string) (string, error)
}

type TokenValidator interface {
	Validate(tokenString string) (*model.Token, error)
}
