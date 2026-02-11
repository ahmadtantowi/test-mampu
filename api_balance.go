package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type BalanceResp struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

func (api *WalletApi) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userId <= 0 {
		sendJSONError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	balance, err := api.getBalance(userId)
	if err != nil {
		api.logger.Error("failed to get user balance", "error", err)
		sendJSONError(w, http.StatusInternalServerError, "failed to get user balance")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func (api *WalletApi) getBalance(ctx context.Context, userId int) (float64, error) {
	const query = `SELECT balance FROM wallets w WHERE w.user_id = $1`

	var balance float64
	if err := api.db.QueryRowxContext(ctx, query, userId).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
