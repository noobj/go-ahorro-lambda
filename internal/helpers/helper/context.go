package helper

import (
	"context"

	"github.com/noobj/go-serverless-services/internal/config"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ContextKeyUser = contextKey("user")

func GetUserFromContext(ctx context.Context) (UserRepository.User, bool) {
	user, ok := ctx.Value(ContextKeyUser).(UserRepository.User)
	return user, ok
}

var ContextKeyConfig = contextKey("config")

func GetConfigFromContext(ctx context.Context) (config.Specification, bool) {
	spec, ok := ctx.Value(ContextKeyConfig).(config.Specification)
	return spec, ok
}
