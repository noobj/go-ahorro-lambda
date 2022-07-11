package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// MiddlewareFunc is a generic middleware example that takes in a HandlerFunc
// and calls the next middleware in the chain.
func MiddlewareFunc(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return next(ctx, request)
	})
}
