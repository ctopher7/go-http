package impl

import (
	"context"
	"database/sql"
	"time"

	"example.com/m/v2/model"
)

func (r *repository) GetUserByEmail(ctx context.Context, email string) (res model.User, err error) {
	query := `
		SELECT
			id, email, password, role
		FROM
			users
		WHERE
			email = $1
	`
	row := r.Db.QueryRowContext(ctx, query, email)

	err = row.Scan(&res.Id, &res.Email, &res.Password, &res.Role)

	return
}

func (r *repository) InsertUser(ctx context.Context, tx *sql.Tx, user model.User) (id int64, err error) {
	query := `
		INSERT INTO
			users(
				email, password,role, created_at,updated_at
			)
		VALUES
			($1,$2,$3,$4,$4)
		RETURNING
			id
	`
	row := tx.QueryRowContext(ctx, query, user.Email, user.Password, user.Role, time.Now())

	err = row.Scan(&id)

	return
}
