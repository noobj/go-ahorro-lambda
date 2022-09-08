package middleware

import (
	"context"

	"github.com/noobj/go-serverless-services/internal/types"
)

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc[T types.ApiRequest, R types.ApiResponse] func(context.Context, T) (R, error)

// MiddlewareFunc is a generic middleware example that takes in a HandlerFunc
// and calls the next middleware in the chain.
func MiddlewareFunc[T types.ApiRequest, R types.ApiResponse](next HandlerFunc[T, R]) HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		return next(ctx, request)
	}
}
