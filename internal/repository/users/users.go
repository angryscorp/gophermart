package users

import (
	"context"
	"fmt"
	"github.com/angryscorp/gophermart/internal/domain/repository"
)

type Users struct {
}

var _ repository.Users = (*Users)(nil)

func New() Users {
	return Users{}
}

func (u Users) CreateUser(ctx context.Context, username, passwordHash string) error {
	fmt.Println("-->>CreateUser", username, passwordHash)
	return nil
}

func (u Users) CheckUser(ctx context.Context, username, passwordHash string) error {
	fmt.Println("-->>CheckUser", username, passwordHash)
	return nil
}
