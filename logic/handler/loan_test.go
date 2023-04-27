package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"example.com/m/v2/constant"
	u "example.com/m/v2/logic/usecase"
	"example.com/m/v2/model"
	"github.com/golang-jwt/jwt/v5"
)

func Test_NewLoan(t *testing.T) {
	ucMock := new(u.MockUsecase)
	rBody := model.NewLoanReq{
		Amount: 10000,
		Terms:  3,
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
				r: httptest.NewRequest("POST", "/loan", &bytes.Buffer{}),
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
					On("DecodeJwt", []*http.Cookie{}).
					Return(jwt.MapClaims{
						"id": float64(1),
					}, nil).
					Once()
				ucMock.
					On("NewLoan", context.Background(), float64(10000), 3, int64(1)).
					Return(nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/loan", &buf),
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

			h.NewLoan(tt.args.w, tt.args.r)
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

func Test_ApproveLoan(t *testing.T) {
	ucMock := new(u.MockUsecase)
	rBody := model.ApproveLoanReq{
		LoanId: 1,
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
				r: httptest.NewRequest("PUT", "/loan/approve", &bytes.Buffer{}),
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
					On("DecodeJwt", []*http.Cookie{}).
					Return(jwt.MapClaims{
						"id":   float64(1),
						"role": constant.AdminRole,
					}, nil).
					Once()
				ucMock.
					On("ApproveLoan", context.Background(), int64(1)).
					Return(nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("PUT", "/loan/approve", &buf),
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

			h.ApproveLoan(tt.args.w, tt.args.r)
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

func Test_PayLoan(t *testing.T) {
	ucMock := new(u.MockUsecase)
	rBody := model.PayLoanReq{
		LoanId: 1,
		Term:   1,
		Amount: 10000,
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
				r: httptest.NewRequest("POST", "/loan/pay", &bytes.Buffer{}),
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
					On("DecodeJwt", []*http.Cookie{}).
					Return(jwt.MapClaims{
						"id":   float64(1),
						"role": constant.AdminRole,
					}, nil).
					Once()
				ucMock.
					On("PayLoan", context.Background(), float64(10000), int64(1), int64(1), int64(1)).
					Return(nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("POST", "/loan/pay", &buf),
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

			h.PayLoan(tt.args.w, tt.args.r)
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

func Test_GetLoan(t *testing.T) {
	ucMock := new(u.MockUsecase)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		name           string
		mock           func()
		args           args
		wantStatusCode int
		wantBody       model.HttpResLoan
		wantHeader     map[string]string
	}{
		{
			name: "err DecodeJwt",
			mock: func() {
				ucMock.
					On("DecodeJwt", []*http.Cookie{}).
					Return(jwt.MapClaims{}, errors.New("err DecodeJwt")).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/loan", nil),
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: model.HttpResLoan{
				Message: "err DecodeJwt",
			},
			wantHeader: map[string]string{
				constant.HttpHeaderSetContent: constant.HttpHeaderAppJson,
			},
		},
		{
			name: "unauthorized",
			mock: func() {
				ucMock.
					On("DecodeJwt", []*http.Cookie{}).
					Return(jwt.MapClaims{
						"id": "asd",
					}, nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/loan", nil),
			},
			wantStatusCode: http.StatusUnauthorized,
			wantBody: model.HttpResLoan{
				Message: "unauthorized",
			},
			wantHeader: map[string]string{
				constant.HttpHeaderSetContent: constant.HttpHeaderAppJson,
			},
		},
		{
			name: "success",
			mock: func() {
				ucMock.
					On("DecodeJwt", []*http.Cookie{}).
					Return(jwt.MapClaims{
						"id":   float64(1),
						"role": constant.AdminRole,
					}, nil).
					Once()
				ucMock.
					On("GetLoan", context.Background(), int64(1)).
					Return([]model.Loan{
						{
							Id: 1,
						},
					}, nil).
					Once()
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/loan", nil),
			},
			wantStatusCode: 200,
			wantBody: model.HttpResLoan{
				Message: "success",
				Data: []model.Loan{
					{
						Id: 1,
					},
				},
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

			h.GetLoan(tt.args.w, tt.args.r)
			if tt.args.w.Result().StatusCode != tt.wantStatusCode {
				t.Errorf("Status code returned, %d, did not match expected code %d", tt.args.w.Result().StatusCode, tt.wantStatusCode)
			}

			var got model.HttpResLoan
			json.NewDecoder(tt.args.w.Body).Decode(&got)
			if !reflect.DeepEqual(got, tt.wantBody) {
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
