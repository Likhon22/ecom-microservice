package cartRepo

import "github.com/redis/go-redis/v9"

type cart struct {
	db *redis.Client
}

type Cart interface{}

func NewRepo(db *redis.Client) Cart {

	return &cart{
		db: db,
	}

}
