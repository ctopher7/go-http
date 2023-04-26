package impl

import (
	"database/sql"

	r "example.com/m/v2/logic/repository"
	"example.com/m/v2/resource"
)

type repository struct {
	Db *sql.DB
}

func New(res *resource.Resource) r.Repository {
	return &repository{
		Db: res.PostgresDb,
	}
}
