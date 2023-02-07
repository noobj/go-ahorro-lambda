package types

import (
	"context"

	"github.com/noobj/jwtmiddleware/types"
)

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc[T ApiRequest, R ApiResponse] func(context.Context, T) (R, error)

type IIvoker[T types.ApiRequest, R types.ApiResponse] interface {
	Invoke(context.Context, T) (R, error)
}
