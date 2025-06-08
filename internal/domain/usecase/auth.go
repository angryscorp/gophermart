package usecase

import "context"

type Auth interface {
	SignUp(ctx context.Context, username, password string) (string, error)
	SignIn(ctx context.Context, username, password string) (string, error)
}
