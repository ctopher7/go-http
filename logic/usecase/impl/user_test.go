package impl

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"example.com/m/v2/constant"
	repo "example.com/m/v2/logic/repository"
	"example.com/m/v2/model"
	"example.com/m/v2/util"
	"github.com/golang-jwt/jwt/v5"
)

func Test_UserLogin(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		email string
		pass  string
	}

	req := args{
		email: "tes",
		pass:  "tes",
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
		want    string
	}{
		{
			name: "fail GetUserByEmail",
			mock: func() {
				repoMock.
					On("GetUserByEmail", context.Background(), "tes").
					Return(model.User{}, errors.New("err GetUserByEmail")).
					Once()
			},
			wantErr: errors.New("err GetUserByEmail"),
			args:    req,
		},
		{
			name: "fail email not registered",
			mock: func() {
				repoMock.
					On("GetUserByEmail", context.Background(), "tes").
					Return(model.User{}, nil).
					Once()
			},
			wantErr: errors.New("email not registered"),
			args:    req,
		},
		{
			name: "fail BcryptComparePassword",
			mock: func() {
				repoMock.
					On("GetUserByEmail", context.Background(), "tes").
					Return(model.User{
						Id:       1,
						Password: "tes",
					}, nil).
					Once()

				repoMock.
					On("BcryptComparePassword", []byte("tes"), []byte("tes")).
					Return(errors.New("err BcryptComparePassword")).
					Once()
			},
			wantErr: errors.New("err BcryptComparePassword"),
			args:    req,
		},
		{
			name: "fail JwtSign",
			mock: func() {
				repoMock.
					On("GetUserByEmail", context.Background(), "tes").
					Return(model.User{
						Id:       1,
						Password: "tes",
						Role:     "tes",
					}, nil).
					Once()

				repoMock.
					On("BcryptComparePassword", []byte("tes"), []byte("tes")).
					Return(nil).
					Once()

				repoMock.
					On("JwtNew", jwt.MapClaims{
						"id":   int64(1),
						"role": "tes",
					}).
					Return(&jwt.Token{}).
					Once()

				repoMock.
					On("JwtSign", &jwt.Token{}).
					Return("", errors.New("err JwtSign")).
					Once()
			},
			wantErr: errors.New("err JwtSign"),
			args:    req,
		},
		{
			name: "success",
			mock: func() {
				repoMock.
					On("GetUserByEmail", context.Background(), "tes").
					Return(model.User{
						Id:       1,
						Password: "tes",
						Role:     "tes",
					}, nil).
					Once()

				repoMock.
					On("BcryptComparePassword", []byte("tes"), []byte("tes")).
					Return(nil).
					Once()

				repoMock.
					On("JwtNew", jwt.MapClaims{
						"id":   int64(1),
						"role": "tes",
					}).
					Return(&jwt.Token{}).
					Once()

				repoMock.
					On("JwtSign", &jwt.Token{}).
					Return("got", nil).
					Once()
			},
			args: req,
			want: "got",
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			got, err := u.UserLogin(context.Background(), tt.args.email, tt.args.pass)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("UserLogin test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserLogin test failed. want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}

func Test_UserRegister(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		user model.User
	}

	req := args{
		user: model.User{
			Email:    "tes",
			Password: "tes",
		},
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
	}{
		{
			name: "fail BcryptGenerateHash",
			mock: func() {
				repoMock.
					On("BcryptGenerateHash", []byte("tes")).
					Return([]byte(""), errors.New("err BcryptGenerateHash")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err BcryptGenerateHash"),
		},
		{
			name: "fail beginTx",
			mock: func() {
				repoMock.
					On("BcryptGenerateHash", []byte("tes")).
					Return([]byte("tes"), nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(nil, errors.New("err beginTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err beginTx"),
		},
		{
			name: "fail InsertUser",
			mock: func() {
				repoMock.
					On("BcryptGenerateHash", []byte("tes")).
					Return([]byte("tes"), nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertUser", context.Background(), &sql.Tx{}, model.User{
						Email:    "tes",
						Password: "tes",
						Role:     constant.CustomerRole,
					}).
					Return(int64(0), errors.New("err InsertUser")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err InsertUser"),
		},
		{
			name: "user not created",
			mock: func() {
				repoMock.
					On("BcryptGenerateHash", []byte("tes")).
					Return([]byte("tes"), nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertUser", context.Background(), &sql.Tx{}, model.User{
						Email:    "tes",
						Password: "tes",
						Role:     constant.CustomerRole,
					}).
					Return(int64(0), nil).
					Once()
			},
			args:    req,
			wantErr: errors.New("failed create user"),
		},
		{
			name: "fail CommitTx",
			mock: func() {
				repoMock.
					On("BcryptGenerateHash", []byte("tes")).
					Return([]byte("tes"), nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertUser", context.Background(), &sql.Tx{}, model.User{
						Email:    "tes",
						Password: "tes",
						Role:     constant.CustomerRole,
					}).
					Return(int64(1), nil).
					Once()

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(errors.New("err CommitTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err CommitTx"),
		},
		{
			name: "success",
			mock: func() {
				repoMock.
					On("BcryptGenerateHash", []byte("tes")).
					Return([]byte("tes"), nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertUser", context.Background(), &sql.Tx{}, model.User{
						Email:    "tes",
						Password: "tes",
						Role:     constant.CustomerRole,
					}).
					Return(int64(1), nil).
					Once()

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(nil).
					Once()
			},
			args: req,
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			err := u.UserRegister(context.Background(), tt.args.user)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("UserRegister test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
		})
	}
}

func Test_DecodeJwt(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		cookies []*http.Cookie
	}

	req := args{
		cookies: []*http.Cookie{
			{
				Name:  "SID",
				Value: "tes",
			},
		},
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
		want    jwt.MapClaims
	}{
		{
			name: "cookie not found",
			args: args{
				cookies: []*http.Cookie{
					{
						Value: "tes",
					},
				},
			},
			wantErr: errors.New("cookie not found"),
		},
		{
			name:    "fail JwtParse",
			args:    req,
			wantErr: errors.New("err JwtParse"),
			mock: func() {
				repoMock.
					On("JwtParse", "tes").
					Return(nil, errors.New("err JwtParse")).
					Once()
			},
		},
		{
			name: "success",
			args: req,
			mock: func() {
				repoMock.
					On("JwtParse", "tes").
					Return(jwt.MapClaims{
						"a": "b",
					}, nil).
					Once()
			},
			want: jwt.MapClaims{
				"a": "b",
			},
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			got, err := u.DecodeJwt(tt.args.cookies)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("DecodeJwt test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeJwt test failed. want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}
