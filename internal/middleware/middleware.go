package middleware

import (
	"context"

	"github.com/noobj/go-serverless-services/internal/types"
)

// MiddlewareFunc is a generic middleware example that takes in a HandlerFunc
// and calls the next middleware in the chain.
func MiddlewareFunc[T types.ApiRequest, R types.ApiResponse](next types.HandlerFunc[T, R]) types.HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		return next(ctx, request)
	}
}
