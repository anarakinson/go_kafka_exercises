package database

import (
	"database/sql"
	"time"

	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер PostgreSQL
)

type DataBase struct {
	connection *sql.DB
}

func NewDatabase(dsn string) (*DataBase, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Log.Info("Connection to database created")
	return &DataBase{
		connection: db,
	}, nil
}

func (db *DataBase) Close() {
	db.connection.Close()
}
