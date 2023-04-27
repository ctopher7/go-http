package impl

import (
	"reflect"
	"testing"

	"example.com/m/v2/config"
	r "example.com/m/v2/logic/repository"
	u "example.com/m/v2/logic/usecase"
)

func Test_New(t *testing.T) {
	type args struct {
		repo *r.MockRepository
		cfg  *config.Config
	}
	tests := []struct {
		name string
		args args
		want u.Usecase
	}{
		{
			name: "success",
			args: args{
				repo: new(r.MockRepository),
				cfg:  &config.Config{},
			},
			want: &usecase{
				repository: new(r.MockRepository),
				cfg:        &config.Config{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.repo, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fail New")
			}
		})
	}
}
