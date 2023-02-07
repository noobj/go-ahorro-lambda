package helper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/golobby/container/v3"
	CategoryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/category"
	EntryRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/entry"
	LoginInfoRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/logininfo"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	. "github.com/noobj/go-serverless-services/internal/repositories/mocks"
)

func BindIocForTesting(mock *MockIRepository, invoker interface{}) error {
	receiverType := reflect.TypeOf(invoker)
	if receiverType == nil || receiverType.Kind() != reflect.Ptr {
		return errors.New("container: invalid invoker type")
	}

	fakeEntry := EntryRepository.EntryRepository{IRepository: mock}
	fakeUser := UserRepository.UserRepository{IRepository: mock}
	fakeLogin := LoginInfoRepository.LoginInfoRepository{IRepository: mock}
	fakeCategory := CategoryRepository.CategoryRepository{IRepository: mock}

	container.SingletonLazy(func() UserRepository.UserRepository {
		return fakeUser
	})

	container.SingletonLazy(func() LoginInfoRepository.LoginInfoRepository {
		return fakeLogin
	})

	container.SingletonLazy(func() CategoryRepository.CategoryRepository {
		return fakeCategory
	})

	container.SingletonLazy(func() EntryRepository.EntryRepository {
		return fakeEntry
	})

	err := container.Fill(invoker)
	if err != nil {
		return fmt.Errorf("couldn't resolve the invoker")
	}

	return nil
}
