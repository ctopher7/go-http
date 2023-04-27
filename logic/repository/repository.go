package io

import (
	"context"
	"database/sql"

	"example.com/m/v2/model"
	"github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (res model.User, err error)
	JwtNew(claim jwt.MapClaims) *jwt.Token
	BcryptComparePassword(hash, password []byte) error
	JwtSign(token *jwt.Token) (string, error)
	InsertUser(ctx context.Context, tx *sql.Tx, user model.User) (id int64, err error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	RollbackTx(tx *sql.Tx) error
	CommitTx(tx *sql.Tx) error
	BcryptGenerateHash(password []byte) ([]byte, error)
	InsertLoan(ctx context.Context, tx *sql.Tx, loan model.Loan) (id int64, err error)
	InsertRepayment(ctx context.Context, tx *sql.Tx, repayment model.Repayment) (id int64, err error)
	JwtParse(token string) (claims jwt.MapClaims, err error)
	UpdateLoan(ctx context.Context, tx *sql.Tx, loan model.Loan) (err error)
}
