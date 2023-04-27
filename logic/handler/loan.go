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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.DecodeJwt(r.Cookies())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	userId, ok := user["id"].(float64)
	if !ok {
		http.Error(w, "invalid token", http.StatusInternalServerError)
		return
	}
	err = h.Usecase.NewLoan(ctx, req.Amount, req.Terms, int64(userId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Usecase.DecodeJwt(r.Cookies())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userRole, ok := user["role"].(string); !ok || userRole != "ADMIN" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	err = h.Usecase.ApproveLoan(ctx, req.LoanId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(model.HttpRes{
		Message: "success",
	})
}
