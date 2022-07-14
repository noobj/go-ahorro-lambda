package main

import (
	"context"

	"github.com/noobj/swim-crowd-lambda-go/internal/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

}

func main() {
	lambda.Start(middleware.Logging(Handler))
}
