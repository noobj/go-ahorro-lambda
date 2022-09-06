package main_test

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/golobby/container/v3"
	main "github.com/noobj/go-serverless-services/cmd/ahorro/refresh"
	"github.com/noobj/go-serverless-services/internal/repositories"
	mockRepo "github.com/noobj/go-serverless-services/internal/repositories/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var fakeObjId, _ = primitive.ObjectIDFromHex("62badc82d420270009a51019")

var _ = Describe("Refresh", func() {
	var fakeRequest events.APIGatewayV2HTTPRequest

	BeforeEach(func() {
		fakeRequest = events.APIGatewayV2HTTPRequest{}

		ctrl := gomock.NewController(GinkgoT())
		m := mockRepo.NewMockIRepository(ctrl)

		container.NamedSingleton("LoginInfoRepo", func() repositories.IRepository {
			return m
		})

		fakeLoginInfoDoc := bson.M{
			"_id":          fakeObjId,
			"user":         fakeObjId,
			"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjoiNjI3MTA2ZDY3YjJmMjVkZGQzZGFmOTY0IiwiZXhwIjoxNjYwNjY3NjA0fQ.xHBGqenagla6eDBVvlglX3W_cMOvI5fBxv-vQNnuTNw",
		}
		fakeSingleResult := mongo.NewSingleResultFromDocument(fakeLoginInfoDoc, nil, nil)

		m.EXPECT().FindOne(context.TODO(), gomock.Any()).Return(fakeSingleResult).MaxTimes(1)
		m.EXPECT().Disconnect().Return(func() {}).MaxTimes(1)
		os.Setenv("REFRESH_TOKEN_SECRET", "420forlife")
		os.Setenv("REFRESH_TOKEN_EXPIRATION_TIME", "3600")
		os.Setenv("ACCESS_TOKEN_EXPIRATION_TIME", "3600")
	})

	Context("when use proper refresh token", func() {
		It("should return accesstoken", func() {
			fakeRequest.Cookies = []string{
				"refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjoiNjJiYWRjODJkNDIwMjcwMDA5YTUxMDE5In0.e5WtrcxJx0w2J6xTLzOSa6TJdR33PN9hdiipazfKmiY",
			}

			res, err := main.Handler(context.TODO(), fakeRequest)

			header := res.Cookies
			Expect(err).To(BeNil())
			Expect(header[0]).Should(ContainSubstring("access_token"))
			Expect(res.StatusCode).To(Equal(200))
		})
	})

	DescribeTable("Login should failed",
		func(token string) {
			tokenString := fmt.Sprintf("refresh_token=%s", token)
			fakeRequest.Cookies = []string{
				tokenString,
			}

			res, err := main.Handler(context.TODO(), fakeRequest)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(401))
		},
		Entry("When with empty token", ""),
		Entry("When with wrong token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjoiNjJiYWRjODJkNDIwMjcwMDA5YTUxMDE5In0.Zs5-TZZSDCM9ykY2vmqNj2qtec02dr34t3Mr0Y5lLWw"),
	)
})
