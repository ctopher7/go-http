package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/m/v2/constant"
	"example.com/m/v2/model"
)

func (h *Handler) NewLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.HttpHeaderSetContent, constant.HttpHeaderAppJson)

	var req model.NewLoanReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	user, err := h.Usecase.DecodeJwt(r.Cookies())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	ctx := context.Background()
	userId, ok := user["id"].(float64)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: "invalid token",
		})
		return
	}
	err = h.Usecase.NewLoan(ctx, req.Amount, req.Terms, int64(userId))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(model.HttpRes{
		Message: "success",
	})
}

func (h *Handler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.HttpHeaderSetContent, constant.HttpHeaderAppJson)

	var req model.ApproveLoanReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	user, err := h.Usecase.DecodeJwt(r.Cookies())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	if userRole, ok := user["role"].(string); !ok || userRole != "ADMIN" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: "unauthorized",
		})
		return
	}
	ctx := context.Background()
	err = h.Usecase.ApproveLoan(ctx, req.LoanId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(model.HttpRes{
		Message: "success",
	})
}

func (h *Handler) PayLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.HttpHeaderSetContent, constant.HttpHeaderAppJson)

	var req model.PayLoanReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	user, err := h.Usecase.DecodeJwt(r.Cookies())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	userId, ok := user["id"].(float64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: "unauthorized",
		})
		return
	}
	ctx := context.Background()
	err = h.Usecase.PayLoan(ctx, req.Amount, req.LoanId, req.Term, int64(userId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(model.HttpRes{
		Message: "success",
	})
}

func (h *Handler) GetLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.HttpHeaderSetContent, constant.HttpHeaderAppJson)

	user, err := h.Usecase.DecodeJwt(r.Cookies())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}

	userId, ok := user["id"].(float64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: "unauthorized",
		})
		return
	}
	ctx := context.Background()
	got, err := h.Usecase.GetLoan(ctx, int64(userId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.HttpRes{
			Message: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(model.HttpResLoan{
		Message: "success",
		Data:    got,
	})
}
