package main_test

import (
	"context"

	"github.com/golang/mock/gomock"
	container "github.com/golobby/container/v3"
	main "github.com/noobj/go-serverless-services/cmd/swimming/crowdscraper"
	"github.com/noobj/go-serverless-services/cmd/swimming/crowdscraper/matchers"
	"github.com/noobj/go-serverless-services/internal/repositories"
	. "github.com/noobj/go-serverless-services/internal/repositories/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crowdscraper", func() {
	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		m := NewMockIRepository(ctrl)

		container.Singleton(func() repositories.IRepository {
			return m
		})

		m.EXPECT().InsertOne(matchers.Regexp(`^20\d\d-[0-1][0-9]-[0-3][0-9] \d{2}:\d{2}`)).MaxTimes(1)
		m.EXPECT().Disconnect().Return(func() {}).MaxTimes(1)
	})

	Context("when handler run normally", func() {
		It("should return status code 200", func() {
			res, err := main.Handler(context.TODO())
			Expect(res.StatusCode).To(Equal(200))
			Expect(err).To(BeNil())
		})
	})
})
