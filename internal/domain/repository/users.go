package repository

import "context"

type Users interface {
	CreateUser(ctx context.Context, username, passwordHash string) error
	CheckUser(ctx context.Context, username, passwordHash string) error
}
