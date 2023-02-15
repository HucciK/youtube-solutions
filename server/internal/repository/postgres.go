package repository

import (
	"database/sql"
	"yt-solutions-server/config"

	_ "github.com/lib/pq"
)

func NewPostgres(cfg config.DBСonfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.ToDataSource())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
