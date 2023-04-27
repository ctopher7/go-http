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

func (r *repository) GetRepaymentByLoanId(ctx context.Context, loanId int64) (res []model.Repayment, err error) {
	query := `
		SELECT
			id, loan_id, minimum_payment, actual_payment, status, due_date
		FROM
			repayments
		WHERE
			loan_id = $1 
		ORDER BY
			due_date ASC
			
	`

	rows, err := r.Db.QueryContext(ctx, query, loanId)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		temp := model.Repayment{}
		err = rows.Scan(&temp.Id, &temp.LoanId, &temp.MinimumPayment, &temp.ActualPayment, &temp.Status, &temp.DueDate)
		if err != nil {
			return
		}
		res = append(res, temp)
	}

	return
}

func (r *repository) UpdateRepayment(ctx context.Context, tx *sql.Tx, repayment model.Repayment) (err error) {

	query := `
		UPDATE
			repayments
		SET
			actual_payment = COALESCE($1, actual_payment), 
			status = COALESCE($2, status), 
			updated_at = $3			
		WHERE
			id = $4 
	`
	_, err = tx.ExecContext(ctx, query, repayment.ActualPayment, repayment.Status, time.Now(), repayment.Id)

	return
}
