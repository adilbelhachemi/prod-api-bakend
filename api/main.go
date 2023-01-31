package main

import (
	"context"
	"log"
	"pratbacknd/internal/server"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
)

var lambdaHandler *fiberadapter.FiberLambda

func init() {
	server, err := server.New(server.Config{Port: "5000"})
	if err != nil {
		log.Fatalf("Could not create server : %s", err)
	}

	lambdaHandler = fiberadapter.New(server.App)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return lambdaHandler.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
