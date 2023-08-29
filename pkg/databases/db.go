package databases

import (
	"log"

	"github.com/k0msak007/kawaii-shop/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(c config.IDbConfig) *sqlx.DB {
	// Connect
	db, err := sqlx.Connect("pgx", c.Url())

	if err != nil {
		log.Fatalf("Connect to database failed: %v\n", err)
	}

	db.DB.SetMaxOpenConns(c.MaxOpenConns())

	return db
}
