package impl

import (
	"context"
	"database/sql"
	"time"

	"example.com/m/v2/model"
)

func (r *repository) InsertRepayment(ctx context.Context, tx *sql.Tx, repayment model.Repayment) (id int64, err error) {
	query := `
		INSERT INTO
			repayments(
				loan_id, minimum_payment, status, due_date, created_at, updated_at
			)
		VALUES
			($1,$2,$3,$4,$5,$5)
		RETURNING
			id
	`
	row := tx.QueryRowContext(ctx, query, repayment.LoanId, repayment.MinimumPayment, repayment.Status, repayment.DueDate, time.Now())

	err = row.Scan(&id)

	return
}
