package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/m/v2/constant"
	u "example.com/m/v2/logic/usecase"
	"example.com/m/v2/model"
)

func Test_UserLogin(t *testing.T) {
	ucMock := new(u.MockUsecase)
	rBody := model.User{
		Email:    "tes@tes.com",
		Password: "tes",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(rBody)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		name           string
		mock           func()
		args           args
		wantStatusCode int
		wantBody       model.HttpRes
		wantHeader     map[string]string
	}{
		{
			name: "err decode req body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/user/login", &bytes.Buffer{}),
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: model.HttpRes{
				Message: "EOF",
			},
			wantHeader: map[string]string{
				constant.HttpHeaderSetContent: constant.HttpHeaderAppJson,
			},
		},
		{
			name: "success",
			mock: func() {
				ucMock.
					On("UserLogin", context.Background(), "tes@tes.com", "tes").
					Return("a", nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/user/login", &buf),
			},
			wantStatusCode: 200,
			wantBody: model.HttpRes{
				Message: "success",
			},
			wantHeader: map[string]string{
				constant.HttpHeaderSetCookie:  fmt.Sprintf("SID=a; HttpOnly; Path=/; Expires=%s; Domain=localhost;", time.Now().Add(24*time.Hour).Format(constant.TimeFormatCookieExpiry)),
				constant.HttpHeaderSetContent: constant.HttpHeaderAppJson,
			},
		},
	}

	for _, tt := range tests {
		h := Handler{
			Usecase: ucMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			h.Login(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("Status code returned, %d, did not match expected code %d", tt.args.w.Result().StatusCode, tt.wantStatusCode)
			}

			var got model.HttpRes
			json.NewDecoder(tt.args.w.Body).Decode(&got)
			if got != tt.wantBody {
				t.Errorf("handler returned unexpected body: got %+v want %+v", got, tt.wantBody)
			}

			for key, val := range tt.wantHeader {
				if tt.args.w.Header().Get(key) != val {
					t.Errorf("handler returned unexpected header: got %+v want %+v", tt.args.w.Header().Get(key), val)
				}
			}
		})
	}
}

func Test_UserRegister(t *testing.T) {
	ucMock := new(u.MockUsecase)
	rBody := model.User{
		Email:    "tes@tes.com",
		Password: "tes",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(rBody)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		name           string
		mock           func()
		args           args
		wantStatusCode int
		wantBody       model.HttpRes
		wantHeader     map[string]string
	}{
		{
			name: "err decode req body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/user/register", &bytes.Buffer{}),
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: model.HttpRes{
				Message: "EOF",
			},
			wantHeader: map[string]string{
				constant.HttpHeaderSetContent: constant.HttpHeaderAppJson,
			},
		},
		{
			name: "success",
			mock: func() {
				ucMock.
					On("UserRegister", context.Background(), model.User{
						Email:    "tes@tes.com",
						Password: "tes",
					}).
					Return(nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/user/register", &buf),
			},
			wantStatusCode: 200,
			wantBody: model.HttpRes{
				Message: "success",
			},
			wantHeader: map[string]string{
				constant.HttpHeaderSetContent: constant.HttpHeaderAppJson,
			},
		},
	}

	for _, tt := range tests {
		h := Handler{
			Usecase: ucMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			h.Register(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("Status code returned, %d, did not match expected code %d", tt.args.w.Result().StatusCode, tt.wantStatusCode)
			}

			var got model.HttpRes
			json.NewDecoder(tt.args.w.Body).Decode(&got)
			if got != tt.wantBody {
				t.Errorf("handler returned unexpected body: got %+v want %+v", got, tt.wantBody)
			}

			for key, val := range tt.wantHeader {
				if tt.args.w.Header().Get(key) != val {
					t.Errorf("handler returned unexpected header: got %+v want %+v", tt.args.w.Header().Get(key), val)
				}
			}
		})
	}
}
