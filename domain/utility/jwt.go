package utility

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	UserID int `json:"UserID"`
	jwt.StandardClaims
}

func NewJWT(userID int, secret string) (string, error) {
	claims := &JwtCustomClaims{
		userID,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
