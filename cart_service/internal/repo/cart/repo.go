package cartRepo

import (
	"cart_service/internal/domain"
	"cart_service/utils"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type repo struct {
	db *redis.Client
}

type Repo interface {
	AddToCart(ctx context.Context, email string, payload *domain.Cart) (*domain.Cart, error)
	GetCart(ctx context.Context, email string) (*domain.Cart, error)
	DeleteCart(ctx context.Context, email string) error
}

func NewRepo(db *redis.Client) Repo {

	return &repo{
		db: db,
	}

}

func (r *repo) AddToCart(ctx context.Context, email string, payload *domain.Cart) (*domain.Cart, error) {
	key := utils.CreateKey(email)
	data, error := json.Marshal(payload)
	if error != nil {
		return nil, error

	}
	_, err := r.db.Set(ctx, key, data, 7*24*time.Hour).Result()
	if err != nil {
		return nil, err

	}

	return payload, nil
}

func (r *repo) GetCart(ctx context.Context, email string) (*domain.Cart, error) {
	key := utils.CreateKey(email)
	val, err := r.db.Get(ctx, key).Result()
	if err == redis.Nil {
		return &domain.Cart{
			Email: email,
			Items: []domain.CartItem{},
		}, nil
	}
	if err != nil {
		return nil, err

	}

	cart := &domain.Cart{}
	if err := json.Unmarshal([]byte(val), cart); err != nil {
		return nil, err
	}

	return cart, nil

}

func (r *repo) DeleteCart(ctx context.Context, email string) error {

	key := utils.CreateKey(email)
	_, err := r.db.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil

}
