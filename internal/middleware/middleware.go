package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
)

type ApiRequest interface {
	events.APIGatewayProxyRequest | helper.APIGatewayV2HTTPRequestWithUser
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
