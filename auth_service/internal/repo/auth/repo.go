package repo

import (
	"auth_service/internal/domain"
	"context"
	"encoding/json"
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
	Get(ctx context.Context, email, deviceID string) (string, time.Time, error)
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

type RefreshTokenData struct {
	Token     string
	ExpiresAt time.Time
}

func (r *authRepo) Store(ctx context.Context, refreshToken, email, device_id string, ttl time.Duration) error {
	now := time.Now().UTC()
	key := r.redisKey(email, device_id)
	data := RefreshTokenData{
		Token:     refreshToken,
		ExpiresAt: now.Add(ttl),
	}

	// Marshal to JSON
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := r.redis.Set(ctx, key, b, ttl).Err(); err != nil {
		return err
	}

	filter := bson.M{"email": email, "device_id": device_id}
	update := bson.M{
		"$set": bson.M{
			"token":      refreshToken,
			"created_at": now,
			"expires_at": data.ExpiresAt,
			"revoked":    false,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = r.authCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err

	}
	return nil

}

func (r *authRepo) Get(ctx context.Context, email, deviceID string) (string, time.Time, error) {

	key := r.redisKey(email, deviceID)

	res, err := r.redis.Get(ctx, key).Result()
	if err == nil {
		var data RefreshTokenData
		if err := json.Unmarshal([]byte(res), &data); err != nil {
			return "", time.Time{}, err
		}
		return data.Token, data.ExpiresAt, nil
	}

	// If Redis miss, fallback to MongoDB
	filter := bson.M{"email": email, "device_id": deviceID}
	var result domain.RefreshTokenDomain
	if err := r.authCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		return "", time.Time{}, err
	}

	data := RefreshTokenData{
		Token:     result.Token,
		ExpiresAt: result.ExpiresAt,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "", time.Time{}, err
	}
	if err := r.redis.Set(ctx, key, b, time.Until(result.ExpiresAt)).Err(); err != nil {
		return "", time.Time{}, err
	}
	return result.Token, result.ExpiresAt, nil
}
