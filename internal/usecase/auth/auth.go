package auth

import (
	"context"
	"github.com/angryscorp/gophermart/internal/domain/model"
	"github.com/angryscorp/gophermart/internal/domain/repository"
	"github.com/angryscorp/gophermart/internal/domain/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const tokenTTL = 24 * time.Hour

type Auth struct {
	repository repository.Users
	jwtSecret  []byte
}

func New(
	repository repository.Users,
	jwtSecret string,
) Auth {
	return Auth{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
	}
}

var _ usecase.Auth = (*Auth)(nil)

func (a Auth) SignUp(ctx context.Context, username, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to hash password")
	}

	id, err := a.repository.CreateUser(ctx, username, string(passwordHash))
	if err != nil {
		return "", errors.Wrap(err, "failed to create user")
	}

	return a.generateToken(id.String())
}

func (a Auth) SignIn(ctx context.Context, username, password string) (string, error) {
	user, err := a.repository.UserData(ctx, username)
	if err != nil {
		return "", errors.Wrap(err, "failed to login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return a.generateToken(user.ID.String())
}

var _ usecase.TokenValidator = (*Auth)(nil)

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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
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
