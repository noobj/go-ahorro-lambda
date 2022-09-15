package jwt_auth_test

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/golobby/container/v3"
	"github.com/joho/godotenv"
	"github.com/noobj/go-serverless-services/internal/helpers/helper"
	jwtAuth "github.com/noobj/go-serverless-services/internal/middleware/jwt_auth"
	"github.com/noobj/go-serverless-services/internal/repositories"
	mockRepo "github.com/noobj/go-serverless-services/internal/repositories/mocks"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var fakeObjId, _ = primitive.ObjectIDFromHex("62badc82d420270009a51019")

var _ = Describe("JwtAuth", func() {
	var fakeRequest events.APIGatewayV2HTTPRequest
	if err := godotenv.Load("../../../.env.example"); err != nil {
		log.Println("No .env file found", err)
	}
	BeforeEach(func() {
		fakeRequest = events.APIGatewayV2HTTPRequest{}
		ctrl := gomock.NewController(GinkgoT())
		m := mockRepo.NewMockIRepository(ctrl)

		container.NamedSingleton("UserRepo", func() repositories.IRepository {
			return m
		})
		os.Setenv("ACCESS_TOKEN_SECRET", "codeeatsleep")

		var fakeUserDoc = bson.M{"_id": fakeObjId, "account": "jjj", "password": "$2b$10$N45EGR5JNu8LlA.VPn5ioe4RxO2XYk0L0PW.vVSxYtS84sBU.Nvye"}
		fakeSingleResult := mongo.NewSingleResultFromDocument(fakeUserDoc, nil, nil)
		m.EXPECT().FindOne(context.TODO(), gomock.Any()).Return(fakeSingleResult).MaxTimes(1)
	})

	Context("when use jwt auth as middleware before handler", func() {
		It("should contains user in context when passing eligible cookie", func() {
			fakeHandler := func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
				user, ok := helper.GetUserFromContext(ctx)
				Expect(ok).To(Equal(true))
				Expect(user.Id).To(Equal(fakeObjId))

				return events.APIGatewayProxyResponse{}, nil
			}
			fakeRequest.Cookies = []string{
				"access_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjoiNjJiYWRjODJkNDIwMjcwMDA5YTUxMDE5In0.K5seedwT5PvKYFeyogMN1F0nurtUfwn2YLc5YFpOLNw",
			}
			jwtAuth.Auth(fakeHandler)(context.Background(), fakeRequest)
		})

		It("should return 404 response when no cookie passed", func() {
			fakeHandler := func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
				return events.APIGatewayProxyResponse{}, nil
			}
			res, err := jwtAuth.Auth(fakeHandler)(context.Background(), fakeRequest)

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(401))
			Expect(res.Body).To(Equal("please login in"))
		})

	})

})
