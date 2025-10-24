package domain

import "time"

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type RefreshTokenDomain struct {
	ID        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email"`
	Token     string    `bson:"token"`
	DeviceID  string    `bson:"device_id"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	Revoked   bool      `bson:"revoked"`
}
