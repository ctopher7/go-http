package impl

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"example.com/m/v2/constant"
	repo "example.com/m/v2/logic/repository"
	"example.com/m/v2/model"
	"example.com/m/v2/util"
	"github.com/stretchr/testify/mock"
)

func Test_NewLoan(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		amount float64
		terms  int
		userId int64
	}

	req := args{
		amount: 10000,
		terms:  3,
		userId: 1,
	}

	reqInsertLoan := model.Loan{
		UserId: &req.userId,
		Amount: &req.amount,
		Status: constant.LoanStatusPending,
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
	}{
		{
			name: "fail beginTx",
			mock: func() {
				repoMock.
					On("BeginTx", context.Background()).
					Return(nil, errors.New("err beginTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err beginTx"),
		},
		{
			name: "fail InsertLoan",
			mock: func() {
				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertLoan", context.Background(), &sql.Tx{}, reqInsertLoan).
					Return(int64(0), errors.New("err InsertLoan")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err InsertLoan"),
		},
		{
			name: "loan not created",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertLoan", context.Background(), &sql.Tx{}, reqInsertLoan).
					Return(int64(0), nil).
					Once()
			},
			args:    req,
			wantErr: errors.New("failed create loan"),
		},
		{
			name: "failed InsertRepayment",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertLoan", context.Background(), &sql.Tx{}, reqInsertLoan).
					Return(int64(1), nil).
					Once()

				repoMock.
					On("InsertRepayment", context.Background(), &sql.Tx{}, mock.Anything).
					Return(int64(0), errors.New("err InsertRepayment")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err InsertRepayment"),
		},
		{
			name: "repayment not created",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertLoan", context.Background(), &sql.Tx{}, reqInsertLoan).
					Return(int64(1), nil).
					Once()

				repoMock.
					On("InsertRepayment", context.Background(), &sql.Tx{}, mock.Anything).
					Return(int64(0), nil).
					Once()
			},
			args:    req,
			wantErr: errors.New("failed create repayment"),
		},
		{
			name: "fail CommitTx",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertLoan", context.Background(), &sql.Tx{}, reqInsertLoan).
					Return(int64(1), nil).
					Once()

				repoMock.
					On("InsertRepayment", context.Background(), &sql.Tx{}, mock.Anything).
					Return(int64(1), nil).
					Times(3)

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(errors.New("err CommitTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err CommitTx"),
		},
		{
			name: "success",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("InsertLoan", context.Background(), &sql.Tx{}, reqInsertLoan).
					Return(int64(1), nil).
					Once()

				repoMock.
					On("InsertRepayment", context.Background(), &sql.Tx{}, mock.Anything).
					Return(int64(1), nil).
					Times(3)

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(nil).
					Once()
			},
			args: req,
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			err := u.NewLoan(context.Background(), tt.args.amount, tt.args.terms, tt.args.userId)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("NewLoan test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
		})
	}
}

func Test_ApproveLoan(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		loanId int64
	}

	req := args{
		loanId: 1,
	}

	reqUpdateLoan := model.Loan{
		Id:     1,
		Status: constant.LoanStatusApproved,
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
	}{
		{
			name: "fail beginTx",
			mock: func() {
				repoMock.
					On("BeginTx", context.Background()).
					Return(nil, errors.New("err beginTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err beginTx"),
		},
		{
			name: "fail UpdateLoan",
			mock: func() {
				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("UpdateLoan", context.Background(), &sql.Tx{}, reqUpdateLoan).
					Return(errors.New("err UpdateLoan")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err UpdateLoan"),
		},
		{
			name: "fail CommitTx",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("UpdateLoan", context.Background(), &sql.Tx{}, reqUpdateLoan).
					Return(nil).
					Once()

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(errors.New("err CommitTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err CommitTx"),
		},
		{
			name: "success",
			mock: func() {

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("UpdateLoan", context.Background(), &sql.Tx{}, reqUpdateLoan).
					Return(nil).
					Once()

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(nil).
					Once()
			},
			args: req,
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			err := u.ApproveLoan(context.Background(), tt.args.loanId)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("ApproveLoan test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
		})
	}
}

func Test_PayLoan(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		loanId int64
		amount float64
		term   int64
		userId int64
	}

	req := args{
		loanId: 1,
		amount: 4000,
		term:   2,
		userId: 1,
	}

	amt := float64(10000)
	getLoanByIdAndUserIdRes := model.Loan{
		Status: constant.LoanStatusApproved,
		Amount: &amt,
	}

	actualPay := float64(3333.33)
	GetRepaymentByLoanIdRes := []model.Repayment{
		{
			Id:             1,
			Status:         constant.RepaymentStatusPaid,
			MinimumPayment: 3333.33,
			ActualPayment:  &actualPay,
		},
		{
			Id:             2,
			Status:         constant.RepaymentStatusPending,
			MinimumPayment: 3333.33,
		},
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
	}{
		{
			name: "fail GetLoanByIdAndUserId",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(model.Loan{}, errors.New("err GetLoanByIdAndUserId")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err GetLoanByIdAndUserId"),
		},
		{
			name: "loan not approved",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(model.Loan{}, nil).
					Once()
			},
			args:    req,
			wantErr: errors.New("loan not approved"),
		},
		{
			name: "err GetRepaymentByLoanId",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(nil, errors.New("err GetRepaymentByLoanId")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err GetRepaymentByLoanId"),
		},
		{
			name: "a term before that has not been paid",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return([]model.Repayment{
						{
							Status: constant.RepaymentStatusPending,
						},
						{
							Status: constant.RepaymentStatusPending,
						},
					}, nil).
					Once()
			},
			args:    req,
			wantErr: errors.New("there is a term before that has not been paid"),
		},
		{
			name: "already paid for this term",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return([]model.Repayment{
						{
							Status: constant.RepaymentStatusPaid,
						},
						{
							Status: constant.RepaymentStatusPaid,
						},
					}, nil).
					Once()
			},
			args:    req,
			wantErr: errors.New("already paid for this term"),
		},
		{
			name: "minimum payment not reached",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				actualPay := 3333.33
				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return([]model.Repayment{
						{
							Status:         constant.RepaymentStatusPaid,
							MinimumPayment: 3333.33,
							ActualPayment:  &actualPay,
						},
						{
							Status:         constant.RepaymentStatusPending,
							MinimumPayment: 3333.33,
						},
					}, nil).
					Once()
			},
			args: args{
				loanId: 1,
				amount: 2000,
				term:   2,
				userId: 1,
			},
			wantErr: errors.New("minimum payment not reached"),
		},
		{
			name: "paid more than loan",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				actualPay := float64(8000)
				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return([]model.Repayment{
						{
							Status:         constant.RepaymentStatusPaid,
							MinimumPayment: 3333.33,
							ActualPayment:  &actualPay,
						},
						{
							Status:         constant.RepaymentStatusPending,
							MinimumPayment: 3333.33,
						},
					}, nil).
					Once()
			},
			args: args{
				loanId: 1,
				amount: 2001,
				term:   2,
				userId: 1,
			},
			wantErr: errors.New("paid more than loan"),
		},
		{
			name: "fail beginTx",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(GetRepaymentByLoanIdRes, nil).
					Once()
				repoMock.
					On("BeginTx", context.Background()).
					Return(nil, errors.New("err beginTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err beginTx"),
		},
		{
			name: "fail UpdateLoan",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(GetRepaymentByLoanIdRes, nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("UpdateLoan", context.Background(), &sql.Tx{}, model.Loan{
						Id:     1,
						Status: constant.LoanStatusPaid,
					}).
					Return(errors.New("err UpdateLoan")).
					Once()
			},
			args: args{
				loanId: 1,
				amount: 6666.67,
				term:   2,
				userId: 1,
			},
			wantErr: errors.New("err UpdateLoan"),
		},
		{
			name: "fail UpdateRepayment for remaining repayment",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return([]model.Repayment{
						{
							Id:             1,
							Status:         constant.RepaymentStatusPaid,
							MinimumPayment: 3333.33,
							ActualPayment:  &actualPay,
						},
						{
							Id:             2,
							Status:         constant.RepaymentStatusPending,
							MinimumPayment: 3333.33,
						},
						{
							Id:             3,
							Status:         constant.RepaymentStatusPending,
							MinimumPayment: 3333.34,
						},
					}, nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				repoMock.
					On("UpdateLoan", context.Background(), &sql.Tx{}, model.Loan{
						Id:     1,
						Status: constant.LoanStatusPaid,
					}).
					Return(nil).
					Once()

				temp := float64(0)
				repoMock.
					On("UpdateRepayment", context.Background(), &sql.Tx{}, model.Repayment{
						Id:            3,
						Status:        constant.RepaymentStatusPaid,
						ActualPayment: &temp,
					}).
					Return(errors.New("err UpdateRepayment")).
					Once()
			},
			args: args{
				loanId: 1,
				amount: 6666.67,
				term:   2,
				userId: 1,
			},
			wantErr: errors.New("err UpdateRepayment"),
		},
		{
			name: "fail UpdateRepayment",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(GetRepaymentByLoanIdRes, nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				temp := float64(4000)
				repoMock.
					On("UpdateRepayment", context.Background(), &sql.Tx{}, model.Repayment{
						Id:            2,
						Status:        constant.RepaymentStatusPaid,
						ActualPayment: &temp,
					}).
					Return(errors.New("err UpdateRepayment")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err UpdateRepayment"),
		},
		{
			name: "fail CommitTx",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(GetRepaymentByLoanIdRes, nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				temp := float64(4000)
				repoMock.
					On("UpdateRepayment", context.Background(), &sql.Tx{}, model.Repayment{
						Id:            2,
						Status:        constant.RepaymentStatusPaid,
						ActualPayment: &temp,
					}).
					Return(nil).
					Once()

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(errors.New("err CommitTx")).
					Once()
			},
			args:    req,
			wantErr: errors.New("err CommitTx"),
		},
		{
			name: "success",
			mock: func() {
				repoMock.
					On("GetLoanByIdAndUserId", context.Background(), int64(1), int64(1)).
					Return(getLoanByIdAndUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(GetRepaymentByLoanIdRes, nil).
					Once()

				repoMock.
					On("BeginTx", context.Background()).
					Return(&sql.Tx{}, nil).
					Once()

				repoMock.
					On("RollbackTx", &sql.Tx{}).
					Return(nil).
					Once()

				temp := float64(4000)
				repoMock.
					On("UpdateRepayment", context.Background(), &sql.Tx{}, model.Repayment{
						Id:            2,
						Status:        constant.RepaymentStatusPaid,
						ActualPayment: &temp,
					}).
					Return(nil).
					Once()

				repoMock.
					On("CommitTx", &sql.Tx{}).
					Return(nil).
					Once()
			},
			args: req,
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			err := u.PayLoan(context.Background(), tt.args.amount, tt.args.loanId, tt.args.term, tt.args.userId)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("PayLoan test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
		})
	}
}

func Test_GetLoan(t *testing.T) {
	repoMock := new(repo.MockRepository)

	type args struct {
		userId int64
	}

	req := args{
		userId: 1,
	}

	getLoanByUserIdRes := []model.Loan{
		{
			Id: 1,
		},
	}

	tests := []struct {
		name    string
		mock    func()
		args    args
		wantErr error
		want    []model.Loan
	}{
		{
			name: "fail GetLoanByUserId",
			mock: func() {
				repoMock.
					On("GetLoanByUserId", context.Background(), int64(1)).
					Return(nil, errors.New("err GetLoanByUserId")).
					Once()
			},
			wantErr: errors.New("err GetLoanByUserId"),
			args:    req,
		},
		{
			name: "fail GetRepaymentByLoanId",
			mock: func() {
				repoMock.
					On("GetLoanByUserId", context.Background(), int64(1)).
					Return(getLoanByUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return(nil, errors.New("err GetRepaymentByLoanId")).
					Once()
			},
			wantErr: errors.New("err GetRepaymentByLoanId"),
			args:    req,
			want: []model.Loan{
				{
					Id: 1,
				},
			},
		},
		{
			name: "success",
			mock: func() {
				repoMock.
					On("GetLoanByUserId", context.Background(), int64(1)).
					Return(getLoanByUserIdRes, nil).
					Once()

				repoMock.
					On("GetRepaymentByLoanId", context.Background(), int64(1)).
					Return([]model.Repayment{
						{
							Id: 1,
						},
					}, nil).
					Once()
			},
			args: req,
			want: []model.Loan{
				{
					Id: 1,
					Repayment: &[]model.Repayment{
						{
							Id: 1,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		u := usecase{
			repository: repoMock,
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			got, err := u.GetLoan(context.Background(), tt.args.userId)
			if !util.SameErrorMessage(err, tt.wantErr) {
				t.Errorf("GetLoan test failed. wantErr: %+v, gotErr: %+v", tt.wantErr, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLoan test failed. want: %+v, got: %+v", tt.want, got)
			}
		})
	}
}
