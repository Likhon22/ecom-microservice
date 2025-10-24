package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaim struct {
	Email string
	Role  string
	jwt.RegisteredClaims
}

func SignedToken(email, role, jwt_secret string, jwt_expire time.Duration) (string, error) {
	if jwt_expire == 0 {
		jwt_expire = 5 * time.Minute
	}
	claims := MyClaim{
		Email: email,
		Role:  role,
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
