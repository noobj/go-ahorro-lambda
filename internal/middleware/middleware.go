package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
)

type ApiRequest interface {
	events.APIGatewayProxyRequest | events.APIGatewayV2HTTPRequest
}

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc[T ApiRequest, R helper.ApiResponse] func(context.Context, T) (R, error)

// MiddlewareFunc is a generic middleware example that takes in a HandlerFunc
// and calls the next middleware in the chain.
func MiddlewareFunc[T ApiRequest, R helper.ApiResponse](next HandlerFunc[T, R]) HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		return next(ctx, request)
	}
}
