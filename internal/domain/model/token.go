package model

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
