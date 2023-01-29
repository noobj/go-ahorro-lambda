package bindioc

import (
	"context"

	"github.com/golobby/container/v3"
	"github.com/noobj/go-serverless-services/internal/repositories"
	"github.com/noobj/jwtmiddleware/types"

	AhorroRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro"
	CategoryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/category"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
)

func Handle[T types.ApiRequest, R types.ApiResponse](next types.HandlerFunc[T, R]) types.HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		userRepo := UserRepository.New()
		defer userRepo.Disconnect()()

		container.NamedSingletonLazy("UserRepo", func() repositories.IRepository {
			return userRepo
		})

		container.NamedSingletonLazy("LoginInfoRepo", func() repositories.IRepository {
			return LoginInfoRepository.New()
		})

		container.NamedSingletonLazy("CategoryRepo", func() repositories.IRepository {
			return CategoryRepository.New()
		})

		container.NamedSingletonLazy("EntryRepo", func() repositories.IRepository {
			return AhorroRepository.New()
		})

		return next(ctx, request)
	}
}
