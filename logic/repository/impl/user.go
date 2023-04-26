package impl

import (
	"context"

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
