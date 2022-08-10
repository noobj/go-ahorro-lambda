package helper

import (
	"context"

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
