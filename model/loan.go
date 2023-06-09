package model

import "time"

type Loan struct {
	Id        int64        `db:"id" json:"id,omitempty"`
	UserId    *int64       `db:"user_id" json:"user_id,omitempty"`
	Amount    *float64     `db:"amount" json:"amount,omitempty"`
	Status    string       `db:"status" json:"status,omitempty"`
	CreatedAt time.Time    `db:"created_at" json:"created_at,omitempty"`
	Repayment *[]Repayment `json:"repayment,omitempty"`
}

type NewLoanReq struct {
	Amount float64 `json:"amount"`
	Terms  int     `json:"terms"`
}

type ApproveLoanReq struct {
	LoanId int64 `json:"loan_id"`
}

type PayLoanReq struct {
	LoanId int64   `json:"loan_id"`
	Term   int64   `json:"term"`
	Amount float64 `json:"amount"`
}
