package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type (
	WithdrawReq struct {
		Amount float64 `json:"amount"`
	}

	WithdrawResp struct {
		UserID  int     `json:"user_id"`
		Balance float64 `json:"balance"`
	}
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func (api *WalletApi) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || userId <= 0 {
		sendJSONError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req WithdrawReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Amount <= 0 {
		sendJSONError(w, http.StatusBadRequest, "invalid amount")
		return
	}

	var balance float64
	if balance, err = api.withdraw(r.Context(), userId, req.Amount); err != nil {
		api.logger.Error("failed to withdraw balance", "error", err)
		if errors.Is(err, ErrInsufficientBalance) {
			sendJSONError(w, http.StatusBadRequest, err.Error())
		} else {
			sendJSONError(w, http.StatusInternalServerError, "failed to withdraw balance")
		}
		return
	}

	sendJSON(w, http.StatusOK, WithdrawResp{
		UserID:  userId,
		Balance: balance,
	})
}

func (api *WalletApi) withdraw(ctx context.Context, userId int, amount float64) (float64, error) {
	tx, err := api.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	const withdrawQuery = `
	UPDATE wallets
	SET
		balance = balance - $2,
		updated_at = now()
	WHERE 1=1
		AND user_id = $1
		AND balance >= $2
	RETURNING balance`

	var balance float64
	if err := tx.QueryRowxContext(ctx, withdrawQuery, userId, amount).Scan(&balance); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrInsufficientBalance
		}
		return 0, err
	}

	const trxQuery = `INSERT INTO transactions (user_id, amount, transaction_type) VALUES ($1, $2, $3)`
	if _, err := tx.ExecContext(ctx, trxQuery, userId, amount, "withdrawal"); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return balance, nil
}
