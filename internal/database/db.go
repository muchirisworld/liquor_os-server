package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func New(cfg *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
    )
	
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open local database: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	
	return db, nil
}