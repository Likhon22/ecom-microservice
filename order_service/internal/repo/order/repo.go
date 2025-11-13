package orderRepo

import (
	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

type Repo interface {
}

func NewRepo(db *sqlx.DB) Repo {

	return &repo{
		db: db,
	}
}
