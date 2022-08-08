package middleware

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noobj/go-serverless-services/internal/middleware"
)

func Logging(f middleware.HandlerFunc[events.APIGatewayProxyRequest]) middleware.HandlerFunc[events.APIGatewayProxyRequest] {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		r = events.APIGatewayProxyRequest(r)
		log.Printf("remote_addr: %s", r.RequestContext.Identity.SourceIP)
		response, err := f(ctx, r)
		return response, err
	}
}
