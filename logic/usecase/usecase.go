package ohlc

import "context"

type Usecase interface {
	UserLogin(ctx context.Context, email, password string) (token string, err error)
}
