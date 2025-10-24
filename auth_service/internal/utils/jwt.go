package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaim struct {
	Email    string
	Role     string
	DeviceId string
	jwt.RegisteredClaims
}

func SignedToken(email, role, device_id, jwt_secret string, jwt_expire time.Duration) (string, error) {
	if jwt_expire == 0 {
		jwt_expire = 5 * time.Minute
	}
	claims := MyClaim{
		Email:    email,
		Role:     role,
		DeviceId: device_id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwt_expire)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwt_secret))
	if err != nil {
		return "", err

	}
	return signedToken, nil

}

func ParseJwt(refreshToken string, jwtSecret string) (*MyClaim, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &MyClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err

	}
	claim, ok := token.Claims.(*MyClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claim, nil
}
