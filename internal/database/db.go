package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func New(cfg *DBConfig) (*sqlx.DB, error) {
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.Database,
	}
	q := dbURL.Query()
	q.Set("sslmode", cfg.SSLMode)
	dbURL.RawQuery = q.Encode()

	db, err := sqlx.Open("pgx", dbURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
