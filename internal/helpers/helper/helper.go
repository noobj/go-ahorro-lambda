package helper

import (
	"bytes"
	"encoding/json"

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
