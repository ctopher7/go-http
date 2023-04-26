package migration

import (
	"fmt"

	"example.com/m/v2/resource"
)

func Migrate(res *resource.Resource) {
	_, err := res.PostgresDb.Query(`
		CREATE TYPE UserRole AS ENUM ('ADMIN','CUSTOMER');
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = res.PostgresDb.Query(`
		CREATE TABLE IF NOT EXISTS users(
			id BIGSERIAL PRIMARY KEY,
			email TEXT UNIQUE,
			password TEXT,
			role UserRole,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		);
	`)
	if err != nil {
		fmt.Println(err)
	}
}

func Seed(res *resource.Resource) {
	_, err := res.PostgresDb.Query(`
		INSERT INTO users(email,password,role,created_at,updated_at) VALUES ('admin@admin.com','$2a$10$DnOPfZCTGIsFTmue/g.wJuaDfr.CCcpYW6y8MqJxnq3AJATTNmRwm','ADMIN',NOW(),NOW())
	`)
	if err != nil {
		fmt.Println(err)
	}
}
