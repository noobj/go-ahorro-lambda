package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type ApiRequest interface {
	events.APIGatewayProxyRequest | events.APIGatewayV2HTTPRequest
}

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc[T ApiRequest] func(context.Context, T) (events.APIGatewayProxyResponse, error)

// MiddlewareFunc is a generic middleware example that takes in a HandlerFunc
// and calls the next middleware in the chain.
func MiddlewareFunc[T ApiRequest](next HandlerFunc[T]) HandlerFunc[T] {
	return func(ctx context.Context, request T) (events.APIGatewayProxyResponse, error) {
		return next(ctx, request)
	}
}
