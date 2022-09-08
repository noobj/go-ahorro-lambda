package types

import "github.com/aws/aws-lambda-go/events"

type ApiResponse interface {
	events.APIGatewayProxyResponse | events.APIGatewayV2HTTPResponse
}

type ApiRequest interface {
	events.APIGatewayProxyRequest | events.APIGatewayV2HTTPRequest
}
