package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)

	db, err := ConnectDB(logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	const port = 8080
	mux := http.NewServeMux()

	logger.Info(fmt.Sprintf("server started on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
