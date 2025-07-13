package store

import (
	"database/sql"
	"go-web/internal/core/ports"
	"go-web/internal/platform"
	"log/slog"

	_ "github.com/lib/pq"
)

type pgStore struct {
	db *sql.DB
}

func NewPg(cfg *platform.Config) ports.Store {
	db, err := sql.Open("postgres", cfg.StoreAddr())
	if err != nil {
		slog.Error("failed open db connection", "error=", err.Error())
		return nil
	}
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping db", "error=", err.Error())
		return nil
	}
	slog.Info("db connected")
	return &pgStore{db: db}
}
