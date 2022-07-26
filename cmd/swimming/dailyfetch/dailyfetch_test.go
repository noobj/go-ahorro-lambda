package main_test

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	container "github.com/golobby/container/v3"
	main "github.com/noobj/swim-crowd-lambda-go/cmd/swimming/dailyfetch"
	"github.com/noobj/swim-crowd-lambda-go/internal/repositories"
	. "github.com/noobj/swim-crowd-lambda-go/internal/repositories/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
)

var _ = Describe("Dailyfetch", func() {
	var fakeData = []bson.M{
		{"Date": "2022-07-13",
			"Entries": []bson.M{
				{
					"Amount": 1234,
					"Time":   "2022-07-13 15:00",
				},
			},
		},
	}
	var fakeRequest events.APIGatewayProxyRequest
	var mockRepo *MockIRepository

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		mockRepo = NewMockIRepository(ctrl)

		container.Singleton(func() repositories.IRepository {
			return mockRepo
		})

		mockRepo.EXPECT().Disconnect().Return(func() {}).MaxTimes(1)
	})

	Context("when handler return expected json response", func() {
		It("normal case", func() {
			expectedRes := "[{\"date\":\"2022-07-13 (Wed)\",\"entries\":[{\"amount\":1234,\"time\":\"2022-07-13 15:00\"}]}]"
			mockRepo.EXPECT().Aggregate(gomock.Any()).Return(fakeData).MaxTimes(1)
			res, err := main.Handler(context.TODO(), fakeRequest)

			Expect(res.Body).To(Equal(expectedRes))
			Expect(err).To(BeNil())
		})

		It("when end greater than start", func() {
			fakeRequest.QueryStringParameters = make(map[string]string)
			fakeRequest.QueryStringParameters["start"] = "2022-01-31"
			fakeRequest.QueryStringParameters["end"] = "2022-01-01"
			expectedRes := "[{\"date\":\"2022-07-13 (Wed)\",\"entries\":[{\"amount\":1234,\"time\":\"2022-07-13 15:00\"}]}]"

			argument := []bson.M{
				{
					"$match": bson.M{
						"$and": bson.A{
							bson.M{"time": bson.M{"$gt": "2022-01-31 00:00:00"}},
							bson.M{"time": bson.M{"$lte": "2022-01-01 23:59:59"}},
						},
					},
				},
				{
					"$group": bson.M{
						"_id": bson.M{
							"$substr": bson.A{"$time", 0, 10},
						},
						"entries": bson.M{
							"$push": bson.M{
								"amount": "$amount",
								"time":   "$time",
							},
						},
					},
				},
			}
			mockRepo.EXPECT().Aggregate(gomock.Eq(argument)).Return(fakeData).MaxTimes(1)
			res, err := main.Handler(context.TODO(), fakeRequest)

			Expect(res.Body).To(Equal(expectedRes))
			Expect(err).To(BeNil())
		})
	})

	Context("when handler will panic", func() {
		It("should panic for wrong query string format", func() {
			var fakeRequest events.APIGatewayProxyRequest
			fakeRequest.QueryStringParameters = make(map[string]string)
			fakeRequest.QueryStringParameters["start"] = "2022"
			fakeRequest.QueryStringParameters["end"] = "2022-31"

			Expect(func() { main.Handler(context.TODO(), fakeRequest) }).Should(PanicWith(MatchRegexp(`^Could not parse time`)))
		})
	})
})
