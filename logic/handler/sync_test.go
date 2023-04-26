package handler

import (
	"errors"
	"testing"

	uc "example.com/m/v2/logic/usecase"
)

func Test_CopyFile(t *testing.T) {
	ucMock := new(uc.MockUsecase)

	type args struct {
		src  string
		dest string
	}

	req := args{
		src:  "tes",
		dest: "tes",
	}

	tests := []struct {
		name string
		mock func()
		args args
	}{
		{
			name: "fail Sync",
			mock: func() {
				ucMock.
					On("Sync", "tes/", "tes/").
					Return(errors.New("err Sync")).
					Once()
			},
			args: req,
		},
		{
			name: "Success req with trailing slash",
			mock: func() {
				ucMock.
					On("Sync", "tes/", "tes/").
					Return(nil).
					Once()
			},
			args: args{
				src:  "tes/",
				dest: "tes/",
			},
		},
		{
			name: "Success",
			mock: func() {
				ucMock.
					On("Sync", "tes/", "tes/").
					Return(nil).
					Once()
			},
			args: req,
		},
	}

	for _, tt := range tests {
		// h := Handler{
		// 	Usecase: ucMock,
		// }

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			// h.NewLoan()
		})
	}
}
