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

	_, err = res.PostgresDb.Query(`
		CREATE TYPE LoanStatus AS ENUM ('PENDING','APPROVED','PAID');
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = res.PostgresDb.Query(`
		CREATE TABLE IF NOT EXISTS loans(
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT,
			amount NUMERIC,
			status LoanStatus,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ
		);
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = res.PostgresDb.Query(`
		CREATE TYPE RepaymentStatus AS ENUM ('PENDING','PAID');
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = res.PostgresDb.Query(`
		CREATE TABLE IF NOT EXISTS repayments(
			id BIGSERIAL PRIMARY KEY,
			loan_id BIGINT,
			minimum_payment NUMERIC,
			actual_payment NUMERIC,
			status RepaymentStatus,
			due_date TIMESTAMPTZ,
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
