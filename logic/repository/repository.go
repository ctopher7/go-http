package io

import (
	"context"

	"example.com/m/v2/model"
	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (res model.User, err error)
	JwtNew(claim jwt.MapClaims) *jwt.Token
	BcryptComparePassword(hash, password []byte) error
	JwtSign(token *jwt.Token) (string, error)
}
