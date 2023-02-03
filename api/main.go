package main

import (
	"context"
	"log"
	"os"
	"pratbacknd/internal/server"
	"pratbacknd/internal/storage"
	"pratbacknd/internal/utils"

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

	storage, err := storage.NewDynamo("ecommerce-dev")
	if err != nil {
		log.Fatalf("Could not create storage interface")
	}

	server, err := server.New(server.Config{Storage: storage, AllowedOrigins: ao, UUIDGen: utils.UUIDV4{}})
	if err != nil {
		log.Fatalf("Could not create server : %s", err)
	}

	chiLambda = chiadapter.New(server.Mux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
