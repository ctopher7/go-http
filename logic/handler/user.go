package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/m/v2/constant"
	"example.com/m/v2/model"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.CookieSetContent, constant.CookieAppJson)

	var req model.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	sid, err := h.Usecase.UserLogin(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(
		constant.CookieSetCookie,
		fmt.Sprintf("SID=%s; HttpOnly; Path=/; Expires=%s; Domain=localhost;", sid, time.Now().Add(24*time.Hour).Format(constant.TimeFormatCookieExpiry)),
	)

	json.NewEncoder(w).Encode(model.HttpRes{
		Message: "success",
	})
}