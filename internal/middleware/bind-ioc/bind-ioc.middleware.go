package bindioc

import (
	"context"

	"github.com/golobby/container/v3"

	// TODO: use internal types instead
	"github.com/noobj/jwtmiddleware/types"

	CategoryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/category"
	EntryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/entry"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
)

func Handle[T types.ApiRequest, R types.ApiResponse](next types.HandlerFunc[T, R]) types.HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		container.SingletonLazy(func() UserRepository.UserRepository {
			return *UserRepository.New()
		})

		container.SingletonLazy(func() LoginInfoRepository.LoginInfoRepository {
			return *LoginInfoRepository.New()
		})

		container.SingletonLazy(func() CategoryRepository.CategoryRepository {
			return *CategoryRepository.New()
		})

		container.SingletonLazy(func() EntryRepository.EntryRepository {
			return *EntryRepository.New()
		})

		return next(ctx, request)
	}
}
