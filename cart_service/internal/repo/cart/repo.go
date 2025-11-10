package cartRepo

import "github.com/redis/go-redis/v9"

type repo struct {
	db *redis.Client
}

type Repo interface{}

func NewRepo(db *redis.Client) Repo {

	return &repo{
		db: db,
	}

}
