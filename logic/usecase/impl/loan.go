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
		UserId: &userId,
		Amount: &amount,
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

func (u *usecase) ApproveLoan(ctx context.Context, loanId int64) (err error) {
	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return
	}
	defer u.repository.RollbackTx(tx)

	err = u.repository.UpdateLoan(ctx, tx, model.Loan{
		Id:     loanId,
		Status: constant.LoanStatusApproved,
	})
	if err != nil {
		return
	}

	err = u.repository.CommitTx(tx)
	return
}

func (u *usecase) PayLoan(ctx context.Context, amount float64, loanId, term, userId int64) (err error) {
	loan, err := u.repository.GetLoanByIdAndUserId(ctx, loanId, userId)
	if err != nil {
		return
	}

	if loan.Status != constant.LoanStatusApproved {
		return errors.New("loan not approved")
	}

	repayments, err := u.repository.GetRepaymentByLoanId(ctx, loanId)
	if err != nil {
		return
	}

	for _, repayment := range repayments[:term-1] {
		if repayment.Status == constant.RepaymentStatusPending {
			return errors.New("there is a term before that has not been paid")
		}
	}

	repaymentData := repayments[term-1]
	if repaymentData.Status == constant.RepaymentStatusPaid {
		return errors.New("already paid for this term")
	}

	paid := float64(0)
	for _, repayment := range repayments {
		if repayment.Status == constant.RepaymentStatusPaid && repayment.ActualPayment != nil {
			paid += *repayment.ActualPayment
		}
	}

	minimumPayment := float64(0)
	for _, repayment := range repayments[0:term] {
		minimumPayment += repayment.MinimumPayment
	}

	if paid+amount < minimumPayment {
		return errors.New("minimum payment not reached")
	}

	if paid+amount > *loan.Amount {
		return errors.New("paid more than loan")
	}

	tx, err := u.repository.BeginTx(ctx)
	if err != nil {
		return
	}
	defer u.repository.RollbackTx(tx)

	if paid+amount == *loan.Amount {
		err = u.repository.UpdateLoan(ctx, tx, model.Loan{
			Id:     loanId,
			Status: constant.LoanStatusPaid,
		})
		if err != nil {
			return
		}

		// release all remaining pending repayment if user already pay before last schedule
		if int(term) < len(repayments) {
			for _, repayment := range repayments[term:] {
				pay := float64(0)
				err = u.repository.UpdateRepayment(ctx, tx, model.Repayment{
					Id:            repayment.Id,
					Status:        constant.RepaymentStatusPaid,
					ActualPayment: &pay,
				})
				if err != nil {
					return
				}
			}
		}
	}

	err = u.repository.UpdateRepayment(ctx, tx, model.Repayment{
		Id:            repaymentData.Id,
		Status:        constant.RepaymentStatusPaid,
		ActualPayment: &amount,
	})
	if err != nil {
		return
	}

	err = u.repository.CommitTx(tx)

	return
}

func (u *usecase) GetLoan(ctx context.Context, userId int64) (loans []model.Loan, err error) {
	loans, err = u.repository.GetLoanByUserId(ctx, userId)
	if err != nil {
		return
	}

	tasks := make(chan model.AsyncTaskRes)
	for idx, loan := range loans {
		go func(loanId int64, index int, out chan<- model.AsyncTaskRes) {
			res, e := u.repository.GetRepaymentByLoanId(ctx, loanId)
			out <- model.AsyncTaskRes{
				Idx: index,
				Err: e,
				Res: res,
			}
		}(loan.Id, idx, tasks)
	}

	for range loans {
		got := <-tasks
		if got.Err != nil {
			err = got.Err
			return
		}

		if val, ok := got.Res.([]model.Repayment); ok {
			loans[got.Idx].Repayment = &val
		}
	}

	return
}
