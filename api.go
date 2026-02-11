package main

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type WalletApi struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewWalletApi(logger *slog.Logger, db *sqlx.DB) *WalletApi {
	return &WalletApi{
		logger: logger,
		db:     db,
	}
}
