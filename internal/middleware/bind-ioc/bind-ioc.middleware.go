package bindioc

import (
	"context"
	"log"
	"reflect"

	"github.com/golobby/container/v3"

	// TODO: use internal types instead
	typesInternal "github.com/noobj/go-serverless-services/internal/types"
	"github.com/noobj/jwtmiddleware/types"

	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	CategoryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/category"
	EntryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/entry"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
)

func Handle[T types.ApiRequest, R types.ApiResponse](invoker typesInternal.IIvoker[T, R]) types.HandlerFunc[T, R] {
	return func(ctx context.Context, request T) (R, error) {
		receiverType := reflect.TypeOf(invoker)
		if receiverType == nil || receiverType.Kind() != reflect.Ptr {
			log.Println("container: invalid invoker type")
			return helper.GenerateErrorResponse[R](500)
		}

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

		err := container.Fill(invoker)
		if err != nil {
			log.Println("Couldn't resolve the invoker")
			return helper.GenerateErrorResponse[R](500)
		}

		return invoker.Invoke(ctx, request)
	}
}
