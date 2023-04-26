package dependency

import (
	"example.com/m/v2/config"
	"example.com/m/v2/logic/handler"
	rImpl "example.com/m/v2/logic/repository/impl"
	"example.com/m/v2/resource"

	ucImpl "example.com/m/v2/logic/usecase/impl"
)

func Init(cfg *config.Config, res *resource.Resource) Dependency {
	repository := rImpl.New(res)
	usecase := ucImpl.New(repository, cfg)

	return Dependency{
		Handler: handler.Handler{
			Usecase: usecase,
		},
	}
}

type Dependency struct {
	Handler handler.Handler
}
