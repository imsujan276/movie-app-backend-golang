package db

import (
	"backend/pkg/config"
	"context"
	"database/sql"
	"time"
)

func OpenDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Db.Driver, cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
