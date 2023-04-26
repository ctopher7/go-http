package model

import "time"

type Repayment struct {
	Id             int64     `db:"id" json:"id,omitempty"`
	LoanId         int64     `db:"loan_id" json:"loan_id,omitempty"`
	MinimumPayment float64   `db:"minimum_payment" json:"minimum_payment,omitempty"`
	ActualPayment  float64   `db:"actual_payment" json:"actual_payment,omitempty"`
	Status         string    `db:"status" json:"status,omitempty"`
	DueDate        time.Time `db:"due_date" json:"due_date,omitempty"`
}
