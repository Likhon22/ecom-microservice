package repo

import (
	"auth_service/internal/domain"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type authRepo struct {
	authCollection *mongo.Collection
	redis          *redis.Client
}
type AuthRepo interface {
	Store(ctx context.Context, refreshToken, email, device_id string, ttl time.Duration) error
}

func (r *authRepo) redisKey(email, deviceID string) string {
	return "refreshToken" + email + ":" + deviceID
}

func NewAuthRepo(db *mongo.Client, redis *redis.Client) AuthRepo {
	collection := db.Database("auth_service").Collection("refresh_tokens")
	return &authRepo{
		authCollection: collection,
		redis:          redis,
	}
}

func (r *authRepo) Store(ctx context.Context, refreshToken, email, device_id string, ttl time.Duration) error {

	key := r.redisKey(email, device_id)

	if err := r.redis.Set(ctx, key, refreshToken, ttl).Err(); err != nil {
		return err
	}

	filter := bson.M{"email": email, "device_id": device_id}
	update := bson.M{
		"$set": bson.M{
			"token":      refreshToken,
			"created_at": time.Now(),
			"expires_at": time.Now().Add(ttl),
			"revoked":    false,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.authCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err

	}
	return nil

}

func (r *authRepo) Get(ctx context.Context, email, device_id string, ttl time.Duration) (string, error) {

	key := r.redisKey(email, device_id)

	token, err := r.redis.Get(ctx, key).Result()
	if err == nil {
		return token, nil

	}
	if err != redis.Nil {
		return "", nil

	}
	filter := bson.M{"email": email, "device_id": device_id}
	var resultResponse *domain.RefreshTokenDomain
	if err := r.authCollection.FindOne(ctx, filter).Decode(&resultResponse); err != nil {
		return "", err
	}
	if err := r.redis.Set(ctx, key, resultResponse.Token, time.Until(resultResponse.ExpiresAt)).Err(); err != nil {
		log.Println("failed to cache token in redis:", err)
	}
	return resultResponse.Token, nil
}
