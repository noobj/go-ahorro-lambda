package main_test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	main "github.com/noobj/go-serverless-services/cmd/ahorro/login"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/golang/mock/gomock"
	container "github.com/golobby/container/v3"
	"github.com/noobj/go-serverless-services/internal/repositories"
	mocks "github.com/noobj/go-serverless-services/internal/repositories/mocks"
)

var fakeObjId, _ = primitive.ObjectIDFromHex("62badc82d420270009a51019")

var fakeUserDoc = bson.M{"_id": fakeObjId, "account": "jjj", "password": "$2b$10$N45EGR5JNu8LlA.VPn5ioe4RxO2XYk0L0PW.vVSxYtS84sBU.Nvye"}

var _ = Describe("Login", func() {
	var fakeRequest events.APIGatewayProxyRequest

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		m := mocks.NewMockIRepository(ctrl)

		container.NamedSingleton("UserRepo", func() repositories.IRepository {
			return m
		})
		container.NamedSingleton("LoginInfoRepo", func() repositories.IRepository {
			return m
		})

		fakeSingleResult := mongo.NewSingleResultFromDocument(fakeUserDoc, nil, nil)

		m.EXPECT().FindOne(context.TODO(), gomock.Any()).Return(fakeSingleResult).MaxTimes(1)
		m.EXPECT().InsertOne(gomock.Any()).Return().MaxTimes(1)
		m.EXPECT().Disconnect().Return(func() {}).MaxTimes(1)
		os.Setenv("ACCESS_TOKEN_EXPIRATION_TIME", "3600")
		os.Setenv("REFRESH_TOKEN_EXPIRATION_TIME", "3600")
	})

	Context("when handler return expected json response", func() {
		It("login should pass", func() {
			var fakeUser = main.LoginDto{
				Account:  "jjj",
				Password: "1234",
			}
			var fakeUserJson, _ = json.Marshal(fakeUser)
			fakeRequest.Body = string(fakeUserJson)
			res, err := main.Handler(context.TODO(), fakeRequest)

			header := res.MultiValueHeaders
			Expect(err).To(BeNil())
			Expect(header["set-cookie"][0]).Should(ContainSubstring("token"))
			Expect(header["set-cookie"][1]).Should(ContainSubstring("token"))
		})
	})

	DescribeTable("Login should failed",
		func(userName string, password string, body string, statusCode int) {
			fakeRequest = events.APIGatewayProxyRequest{}
			if userName != "" {
				var fakeUser = main.LoginDto{
					Account:  userName,
					Password: password,
				}
				var fakeUserJson, _ = json.Marshal(fakeUser)
				fakeRequest.Body = string(fakeUserJson)
			}
			res, err := main.Handler(context.TODO(), fakeRequest)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(statusCode))
			Expect(res.Body).Should(Equal(body))
		},
		Entry("When wrong password", "jjj", "12342", "account and password not match", 404),
		Entry("When no user found", "", "", "request body error", 400),
	)
})
