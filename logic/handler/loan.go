package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/m/v2/model"
)

func (h *Handler) NewLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
