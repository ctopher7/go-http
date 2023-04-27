package impl

import (
	"example.com/m/v2/constant"
	"github.com/golang-jwt/jwt/v5"
)

func (r *repository) JwtNew(claim jwt.MapClaims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
}

func (r *repository) JwtSign(token *jwt.Token) (string, error) {
	return token.SignedString([]byte(constant.JwtSecret))
}

func (r *repository) JwtParse(token string) (claims jwt.MapClaims, err error) {
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.JwtSecret), nil
	})

	return
}
