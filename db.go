package main

import (
	"context"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const schema = `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE wallets (
    user_id SERIAL PRIMARY KEY REFERENCES users(id),
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal');

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    amount DECIMAL(15, 2) NOT NULL,
    transaction_type transaction_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO users (id, username) VALUES
	(1, 'asep'),
	(2, 'bambang');

INSERT INTO wallets (user_id, balance) VALUES
	(1, 1000000),
	(2, 500000);
`

const (
	connstr = "postgres://postgres:abcde@localhost:5432/test_mampu?sslmode=disable"
	driver  = "pgx"
)

func ConnectDB(ctx context.Context, logger *slog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, driver, connstr)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("failed to ping database", "error", err)
		return nil, err
	}

	logger.Info("applying database schema...")
	res, err := db.ExecContext(ctx, schema)
	if err != nil {
		logger.Warn("failed to apply database schema", "error", err)
	} else {
		logger.Info("database schema applied", "rows_affected", res)
	}

	return db, nil
}
