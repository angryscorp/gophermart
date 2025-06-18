package auth

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

type Auth struct {
	repository repository.Users
	jwtSecret  []byte
}

var _ usecase.Auth = (*Auth)(nil)
var _ usecase.TokenValidator = (*Auth)(nil)

func New(
	repository repository.Users,
	jwtSecret string,
) Auth {
	return Auth{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
	}
}

func (a Auth) SignUp(ctx context.Context, username, password string) (string, error) {
	passwordHash := password // TODO
	id, err := a.repository.CreateUser(ctx, username, passwordHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to create user")
	}

	token, err := a.generateToken(id.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return token, nil

}

func (a Auth) SignIn(ctx context.Context, username, password string) (string, error) {
	passwordHash := password // TODO
	id, err := a.repository.CheckUser(ctx, username, passwordHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to login")
	}

	token, err := a.generateToken(id.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return token, nil

}

func (a Auth) Validate(tokenString string) (*model.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Token{}, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if claims, ok := token.Claims.(*model.Token); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (a Auth) generateToken(userID string) (string, error) {
	claims := model.Token{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate token")
	}

	return tokenString, nil
}
