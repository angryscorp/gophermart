package repository

import (
	"context"
	"github.com/google/uuid"
)

type Users interface {
	CreateUser(ctx context.Context, username, passwordHash string) (*uuid.UUID, error)
	CheckUser(ctx context.Context, username, passwordHash string) (*uuid.UUID, error)
}
