package main_test

import (
	"bytes"
	"context"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
	main "github.com/noobj/go-serverless-services/cmd/ahorro/login"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/golang/mock/gomock"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	mocks "github.com/noobj/go-serverless-services/internal/repositories/mocks"
)

var fakeObjId, _ = primitive.ObjectIDFromHex("62badc82d420270009a51019")

var fakeUserDoc = bson.M{"_id": fakeObjId, "account": "jjj", "password": "$2b$10$N45EGR5JNu8LlA.VPn5ioe4RxO2XYk0L0PW.vVSxYtS84sBU.Nvye"}

var _ = Describe("Login", func() {
	var fakeRequest events.APIGatewayProxyRequest
	invoker := main.Invoker{}

	if err := godotenv.Load("../../../.env.example"); err != nil {
		log.Println("No .env file found", err)
	}
	BeforeEach(func() {
		fakeRequest = events.APIGatewayProxyRequest{
			Headers: make(map[string]string),
		}

		ctrl := gomock.NewController(GinkgoT())
		m := mocks.NewMockIRepository(ctrl)
		if err := helper.BindIocForTesting(m, &invoker); err != nil {
			panic(err.Error())
		}

		fakeSingleResult := mongo.NewSingleResultFromDocument(fakeUserDoc, nil, nil)
		fakeInsertOneResult := &mongo.InsertOneResult{InsertedID: 123}

		m.EXPECT().FindOne(context.TODO(), gomock.Any()).Return(fakeSingleResult).MaxTimes(1)
		m.EXPECT().InsertOne(context.TODO(), gomock.Any()).Return(fakeInsertOneResult, nil).MaxTimes(1)
		m.EXPECT().Disconnect().Return(func() {}).MaxTimes(1)
		os.Setenv("ACCESS_TOKEN_EXPIRATION_TIME", "3600")
		os.Setenv("REFRESH_TOKEN_EXPIRATION_TIME", "3600")
	})

	Context("when handler return expected json response", func() {
		It("login should pass", func() {
			var buf bytes.Buffer
			mw := multipart.NewWriter(&buf)
			fakeRequest.Headers["content-type"] = mw.FormDataContentType()
			mw.WriteField("account", "jjj")
			mw.WriteField("password", "1234")
			mw.Close()
			fakeRequest.Body = buf.String()

			res, err := invoker.Invoke(context.TODO(), fakeRequest)

			header := res.Cookies
			Expect(err).To(BeNil())
			Expect(header[0]).Should(ContainSubstring("token"))
			Expect(header[1]).Should(ContainSubstring("token"))
		})
	})

	DescribeTable("Login should failed",
		func(userName string, password string, body string, statusCode int) {
			var buf bytes.Buffer
			mw := multipart.NewWriter(&buf)
			fakeRequest.Headers["content-type"] = mw.FormDataContentType()
			mw.WriteField("account", userName)
			mw.WriteField("password", password)
			mw.Close()
			fakeRequest.Body = buf.String()

			res, err := invoker.Invoke(context.TODO(), fakeRequest)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(statusCode))
			Expect(res.Body).Should(Equal(body))
		},
		Entry("When wrong password", "jjj", "12342", "account and password not match", 404),
		Entry("When no user found", "", "", "request body error", 400),
	)
})
