package ohlc

import (
	"context"
	"net/http"

	"example.com/m/v2/model"
	"github.com/golang-jwt/jwt/v5"
)

type Usecase interface {
	UserLogin(ctx context.Context, email, password string) (token string, err error)
	UserRegister(ctx context.Context, user model.User) (err error)
	NewLoan(ctx context.Context, amount float64, terms int, userId int64) (err error)
	DecodeJwt(cookies []*http.Cookie) (claims jwt.MapClaims, err error)
}
