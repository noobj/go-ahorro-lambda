package middleware

import (
	"context"
	"fmt"

	typesInternal "github.com/noobj/go-serverless-services/internal/types"
	"github.com/noobj/jwtmiddleware/types"
)

type middlewaresFunc[T types.ApiRequest, R types.ApiResponse] func(invoker typesInternal.IIvoker[T, R]) types.HandlerFunc[T, R]

func fakeMiddlewareFunc[T types.ApiRequest, R types.ApiResponse](invoker typesInternal.IIvoker[T, R]) types.HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		fmt.Println("Fake middleware")
		return invoker.Invoke(ctx, request)
	}
}

func Bootstrap[T types.ApiRequest, R types.ApiResponse](invoker typesInternal.IIvoker[T, R], middlewareFuncs ...middlewaresFunc[T, R]) types.HandlerFunc[T, R] {
	result := func(mFunc middlewaresFunc[T, R]) middlewaresFunc[T, R] {
		return func(invoker typesInternal.IIvoker[T, R]) types.HandlerFunc[T, R] {
			return mFunc(invoker)
		}
	}

	for _, middlewareFunc := range middlewareFuncs {
		oldResult := result
		result = func(mFunc middlewaresFunc[T, R]) middlewaresFunc[T, R] {
			return oldResult(middlewareFunc)
		}
	}

	return result(fakeMiddlewareFunc[T, R])(invoker)
}
