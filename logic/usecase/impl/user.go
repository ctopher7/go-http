package impl

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func (u *usecase) UserLogin(ctx context.Context, email, password string) (token string, err error) {
	user, err := u.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return
	}
	if user.Id <= 0 {
		err = errors.New("email not registered")
		return
	}

	err = u.repository.BcryptComparePassword([]byte(user.Password), []byte(password))
	if err != nil {
		return
	}

	tkn := u.repository.JwtNew(jwt.MapClaims{
		"id":   user.Id,
		"role": user.Role,
	})

	token, err = u.repository.JwtSign(tkn)

	return
}
