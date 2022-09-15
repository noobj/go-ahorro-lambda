package helper

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/noobj/go-serverless-services/internal/config"
	"github.com/noobj/go-serverless-services/internal/types"
)

func GenerateJwtToken(payload interface{}, expiredTime int64, secret string) (string, error) {
	expiresAt := time.Now().Add(time.Duration(expiredTime) * time.Second).Unix()
	claims := types.MyCustomClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func GenerateAccessToken(userId string) (string, error) {
	env := config.GetInstance()

	token, err := GenerateJwtToken(userId, env.AccessTokenExpirationTime, env.AccessTokenSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateRefreshToken(userId string) (string, error) {
	env := config.GetInstance()
	token, err := GenerateJwtToken(userId, env.RefreshTokenExpirationTime, env.RefreshTokenSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
