package main

import (
	"context"
	"log"
	"os"
	"pratbacknd/internal/server"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
)

// var lambdaHandler *chiadapter.
var chiLambda *chiadapter.ChiLambda

func init() {
	ao, found := os.LookupEnv("ALLOWED_ORIGIN")
	if !found {
		log.Fatalf("Could not find ALLOWED_ORIGIN")
	}

	server, err := server.New(server.Config{Port: "5000", AllowedOrigins: ao})
	if err != nil {
		log.Fatalf("Could not create server : %s", err)
	}

	chiLambda = chiadapter.New(server.Router)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
