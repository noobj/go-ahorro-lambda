package jwt_auth

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noobj/go-serverless-services/internal/middleware"
)

func Auth(f middleware.HandlerFunc[events.APIGatewayV2HTTPRequest]) middleware.HandlerFunc[events.APIGatewayV2HTTPRequest] {
	return func(ctx context.Context, r events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
		fmt.Printf("%+v", r.Cookies)

		return f(ctx, r)
	}
}
