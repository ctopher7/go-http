package impl

import (
	"context"
	"errors"
	"reflect"
	"testing"

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

			got, err := u.UserLogin(context.Background(), req.email, req.pass)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("UserLogin test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserLogin test failed. want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}
