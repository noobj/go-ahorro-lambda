package main_test

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/golobby/container/v3"
	main "github.com/noobj/go-serverless-services/cmd/ahorro/fetchentries"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	"github.com/noobj/go-serverless-services/internal/repositories"
	UserRepository "github.com/noobj/go-serverless-services/internal/repositories/ahorro/user"
	. "github.com/noobj/go-serverless-services/internal/repositories/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var fakeObjId, _ = primitive.ObjectIDFromHex("62badc82d420270009a51019")

var fakeData = []bson.M{
	{
		"sum": 110,
		"_id": fakeObjId,
		"category": []bson.M{
			{
				"_id":   fakeObjId,
				"color": "#a4e56c",
				"name":  "Food",
				"user":  fakeObjId,
			},
		},
		"entries": []bson.M{
			{
				"_id":    fakeObjId,
				"amount": 110,
				"date":   "2022-01-05",
				"descr":  "fuck",
			},
		},
	},
	{
		"sum": 90,
		"_id": fakeObjId,
		"category": []bson.M{
			{
				"_id":   fakeObjId,
				"color": "#a4e51c",
				"name":  "Abc",
				"user":  fakeObjId,
			},
		},
		"entries": []bson.M{
			{
				"_id":    fakeObjId,
				"amount": 90,
				"date":   "2022-01-05",
				"descr":  "fuck",
			},
		},
	},
}

var _ = Describe("Fetchentries", func() {
	var fakeRequest events.APIGatewayV2HTTPRequest
	ctx := context.WithValue(context.Background(), helper.ContextKeyUser, UserRepository.User{
		Id:       fakeObjId,
		Account:  "jjj",
		Password: "123456",
	})

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		m := NewMockIRepository(ctrl)

		container.Singleton(func() repositories.IRepository {
			return m
		})
		fakeRequest.QueryStringParameters = make(map[string]string)
		fakeRequest.QueryStringParameters["timeStart"] = "2022-01-01"
		fakeRequest.QueryStringParameters["timeEnd"] = "2022-01-31"

		m.EXPECT().Aggregate(gomock.Any()).Return(fakeData).MaxTimes(1)
		m.EXPECT().Disconnect().Return(func() {}).MaxTimes(1)
	})

	Context("when handler return expected json response", func() {
		It("should pass", func() {
			expectedRes := "{\"categories\":[{\"_id\":\"62badc82d420270009a51019\",\"sum\":110,\"percentage\":\"0.55\",\"name\":\"Food\",\"entries\":[{\"_id\":\"62badc82d420270009a51019\",\"amount\":110,\"date\":\"2022-01-05\",\"descr\":\"fuck\"}],\"color\":\"#a4e56c\"},{\"_id\":\"62badc82d420270009a51019\",\"sum\":90,\"percentage\":\"0.45\",\"name\":\"Abc\",\"entries\":[{\"_id\":\"62badc82d420270009a51019\",\"amount\":90,\"date\":\"2022-01-05\",\"descr\":\"fuck\"}],\"color\":\"#a4e51c\"}],\"total\":200}"
			res, err := main.Handler(ctx, fakeRequest)

			Expect(res.Body).To(Equal(expectedRes))
			Expect(err).To(BeNil())
		})

		It("should panic for wrong query string format", func() {
			res, err := main.Handler(ctx, events.APIGatewayV2HTTPRequest{})
			Expect(res.Body).To(Equal("request query error"))
			Expect(res.StatusCode).To(Equal(400))
			Expect(err).To(BeNil())
		})
	})
})
