package repo

import (
	"auth_service/internal/domain"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type authRepo struct {
	authCollection *mongo.Collection
	redis          *redis.Client
}
type AuthRepo interface {
	Store(ctx context.Context, refreshToken, email string, ttl time.Duration) error
}

func NewAuthRepo(db *mongo.Client, redis *redis.Client) AuthRepo {
	collection := db.Database("auth_service").Collection("refresh_tokens")
	return &authRepo{
		authCollection: collection,
		redis:          redis,
	}
}

func (r *authRepo) Store(ctx context.Context, refreshToken, email string, ttl time.Duration) error {
	if err := r.redis.Set(ctx, "refresh:"+email, refreshToken, ttl); err != nil {
		return err.Err()
	}

	rt := domain.RefreshTokenDomain{
		Email:     email,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		Revoked:   false,
	}
	_, err := r.authCollection.InsertOne(ctx, rt)
	return err

}
