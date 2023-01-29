package types

import (
	"context"
)

// HandlerFunc is a generic JSON Lambda handler used to chain middleware.
type HandlerFunc[T ApiRequest, R ApiResponse] func(context.Context, T) (R, error)
