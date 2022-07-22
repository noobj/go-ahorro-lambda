package helpers

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
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
