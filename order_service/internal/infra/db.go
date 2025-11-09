package infra

import (
	"order_service/internal/config"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	dbInstance *sqlx.DB
	once       sync.Once
)

func ConnectDb(cnf *config.DBConfig) (*sqlx.DB, error) {

	var err error
	once.Do(func() {
		dbInstance, err = sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
		if err != nil {
			return

		}
		err = dbInstance.Ping()
		if err != nil {
			return

		}
	})

	dbInstance.SetMaxOpenConns(cnf.MaxOpenConns)
	dbInstance.SetMaxIdleConns(cnf.MaxIdleConns)
	dbInstance.SetConnMaxLifetime(5 * time.Minute)

	return dbInstance, nil
}
