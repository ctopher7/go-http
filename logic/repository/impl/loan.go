package impl

import (
	"context"
	"database/sql"
	"time"

	"example.com/m/v2/model"
)

func (r *repository) InsertLoan(ctx context.Context, tx *sql.Tx, loan model.Loan) (id int64, err error) {
	query := `
		INSERT INTO
			loans(
				amount, status, user_id, created_at,updated_at
			)
		VALUES
			($1,$2,$3,$4,$4)
		RETURNING
			id
	`
	row := tx.QueryRowContext(ctx, query, loan.Amount, loan.Status, loan.UserId, time.Now())

	err = row.Scan(&id)

	return
}
