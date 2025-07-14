package store

import (
	"database/sql"
	"go-web/internal/core/ports"
	"log/slog"

	_ "github.com/lib/pq"
)

type pgStore struct {
	db *sql.DB
}

func NewPgStore(addr string) ports.Store {
	db, err := sql.Open("postgres", addr)
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
