package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
	"github.com/noobj/go-serverless-services/internal/types"
)

func GenerateApiResponse(resultForBody interface{}) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(resultForBody)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"set-cookie":   "xxx=123",
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func GenerateInternalErrorResponse() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: "internal error", StatusCode: 500}, nil
}

func GenerateJwtToken(payload interface{}, expiredTime int, secret string) (string, error) {
	claims := types.MyCustomClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(expiredTime),
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

func SetCookie(cookie http.Cookie, reps *events.APIGatewayProxyResponse) {
	reps.MultiValueHeaders["set-cookie"] = append(reps.MultiValueHeaders["set-cookie"], cookie.String())
}
