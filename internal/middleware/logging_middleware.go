package middleware

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func Logging(f HandlerFunc) HandlerFunc {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Printf("remote_addr: %s", r.RequestContext.Identity.SourceIP)
		response, err := f(ctx, r)
		return response, err
	}
}
