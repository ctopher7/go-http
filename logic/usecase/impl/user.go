package impl

import (
	"context"
	"errors"
	"net/http"

	"example.com/m/v2/constant"
	"example.com/m/v2/model"
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

func (u *usecase) UserRegister(ctx context.Context, user model.User) (err error) {
	hashPass, err := u.repository.BcryptGenerateHash([]byte(user.Password))
	if err != nil {
		return
	}
	user.Password = string(hashPass)
	user.Role = constant.CustomerRole

	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return
	}
	defer u.repository.RollbackTx(tx)

	id, err := u.repository.InsertUser(ctx, tx, user)
	if err != nil {
		return
	}
	if id <= 0 {
		err = errors.New("failed create user")
		return
	}

	err = u.repository.CommitTx(tx)

	return
}

func (u *usecase) DecodeJwt(cookies []*http.Cookie) (claims jwt.MapClaims, err error) {
	tokenStr := ""
	for _, c := range cookies {
		if c.Name == "SID" {
			tokenStr = c.Value
		}
	}
	if tokenStr == "" {
		err = errors.New("cookie not found")
		return
	}

	_, err = jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.JwtSecret), nil
	})

	return
}
