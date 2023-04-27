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

func (r *repository) UpdateLoan(ctx context.Context, tx *sql.Tx, loan model.Loan) (err error) {
	query := `
		UPDATE
			loans
		SET
			amount = COALESCE($1,amount), 
			status = COALESCE($2,status), 
			user_id = COALESCE($3,user_id),
			updated_at = $4
		WHERE
			id = $5
	`

	_, err = tx.ExecContext(ctx, query, loan.Amount, loan.Status, loan.UserId, time.Now(), loan.Id)

	return
}
