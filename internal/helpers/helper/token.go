package helper

import (
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
	"github.com/noobj/go-serverless-services/internal/types"
)

func GenerateInternalErrorResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: "internal error", StatusCode: 500}, nil
}

func GenerateJwtToken(payload interface{}, expiredTime int, secret string) (string, error) {
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
	accessTokenExpireTime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION_TIME"))
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if err != nil {
		return "", err
	}

	token, err := GenerateJwtToken(userId, accessTokenExpireTime, accessTokenSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GenerateRefreshToken(userId string) (string, error) {
	refreshTokenExpireTime, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION_TIME"))
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if err != nil {
		return "", err
	}

	token, err := GenerateJwtToken(userId, refreshTokenExpireTime, refreshTokenSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
