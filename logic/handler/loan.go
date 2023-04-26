package handler

import (
	"net/http"
)

func (h *Handler) NewLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode()
}
