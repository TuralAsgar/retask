package main

import (
	"context"
	"database/sql"
	"github.com/TuralAsgar/dynamic-programming/internal/data"
	"log/slog"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}

	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// since we are using sqlite
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
