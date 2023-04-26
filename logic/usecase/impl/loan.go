package impl

import (
	"context"
	"errors"
	"math"
	"time"

	"example.com/m/v2/constant"
	"example.com/m/v2/model"
)

func (u *usecase) NewLoan(ctx context.Context, amount float64, terms int, userId int64) (err error) {
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return
	}
	defer u.repository.RollbackTx(tx)

	loanId, err := u.repository.InsertLoan(ctx, tx, model.Loan{
		UserId: userId,
		Amount: amount,
		Status: constant.LoanStatusPending,
	})
	if err != nil {
		return
	}
	if loanId <= 0 {
		err = errors.New("failed create loan")
		return
	}

	minimumPaymentDivided := math.Round(amount/float64(terms)*100) / 100
	repaymentLeft := amount
	for i := 1; i <= terms; i++ {

		minimumPayment := float64(0)
		if i == terms {
			minimumPayment = repaymentLeft
		} else {
			minimumPayment = minimumPaymentDivided
			repaymentLeft -= minimumPaymentDivided
		}

		idRepayment, errRepayment := u.repository.InsertRepayment(ctx, tx, model.Repayment{
			LoanId:         loanId,
			MinimumPayment: minimumPayment,
			Status:         constant.RepaymentStatusPending,
			DueDate:        time.Now().AddDate(0, 0, 7*i),
		})
		if errRepayment != nil {
			err = errRepayment
			return
		}
		if idRepayment <= 0 {
			err = errors.New("failed create repayment")
			return
		}

	}

	err = u.repository.CommitTx(tx)
	return
}
