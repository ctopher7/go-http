package impl

import (
	"example.com/m/v2/config"
	r "example.com/m/v2/logic/repository"
	u "example.com/m/v2/logic/usecase"
)

type usecase struct {
	repository r.Repository
	cfg        *config.Config
}

func New(
	repository r.Repository,
	cfg *config.Config,

) u.Usecase {
	return &usecase{
		repository: repository,
		cfg:        cfg,
	}
}
