package repository

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/google/uuid"
)

type Users interface {
	CreateUser(ctx context.Context, username, passwordHash string) (*uuid.UUID, error)
	UserData(ctx context.Context, username string) (*model.User, error)
}
