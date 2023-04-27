package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	req := args{
		w: httptest.NewRecorder(),
		r: httptest.NewRequest("POST", "/user/login", &buf),
	}

	tests := []struct {
		name           string
		mock           func()
		args           args
		wantStatusCode int
		wantBody       model.HttpRes
	}{
		{
			name: "success",
			mock: func() {
				ucMock.
					On("UserLogin", context.Background(), "tes@tes.com", "tes").
					Return("a", nil).
					Once()
			},
			args:           req,
			wantStatusCode: 200,
			wantBody: model.HttpRes{
				Message: "success",
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
		})
	}
}
