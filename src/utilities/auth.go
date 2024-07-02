package utilities

import (
	"e-learning/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	Status    int    `json:"status"`
	Type      string `json:"type"`
}

func CreateToken(accountID string, username string, role string, status int, expiresTime int, typeToken string) (string, error) {
	claims := MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        primitive.NewObjectID().Hex(),
			Issuer:    accountID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expiresTime)),
		},
		AccountID: accountID,
		Username:  username,
		Role:      role,
		Status:    status,
		Type:      typeToken,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Get().JwtSecret))
	return token, err
}
