package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"pratbacknd/internal/secret"
	"pratbacknd/internal/server"
	"pratbacknd/internal/storage"
	"pratbacknd/internal/utils"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"google.golang.org/api/option"
)

// var lambdaHandler *chiadapter.
var chiLambda *chiadapter.ChiLambda

func init() {
	// allow origins
	ao, found := os.LookupEnv("ALLOWED_ORIGIN")
	if !found {
		log.Fatalf("Could not find ALLOWED_ORIGIN")
	}

	// set firebase
	app, err := setupFireBaseApp()
	if err != nil {
		log.Fatalf("Could not create firebase app : %s", err)
	}

	// set authClient
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Could not build auth client : %s", err)
	}

	// set db
	storage, err := storage.NewDynamo("ecommerce-dev")
	if err != nil {
		log.Fatalf("Could not create storage interface")
	}

	server, err := server.New(
		server.Config{
			Storage:            storage,
			AllowedOrigins:     ao,
			UUIDGen:            utils.UUIDV4{},
			FirebaseAuthClient: authClient,
		},
	)
	if err != nil {
		log.Fatalf("Could not create server : %s", err)
	}

	chiLambda = chiadapter.New(server.Mux)
}

func setupFireBaseApp() (*firebase.App, error) {
	parameterStoreName, found := os.LookupEnv("PARAMETER_STORE_NAME")
	if !found {
		log.Fatalf("Could not find PARAMETER_STORE_NAME")
	}

	ssmClient := ssm.New(session.Must(session.NewSession()))
	outSSM, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameterStoreName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("Could not access SSM : %s", err)
	}

	parameterRawValue := *outSSM.Parameter.Value
	var secrets secret.Parameters
	err = json.Unmarshal([]byte(parameterRawValue), &secrets)
	if err != nil {
		log.Fatalf("Could ont unmarshall secrets: %s", err)
	}

	jsonCreds, err := json.Marshal(secrets.Google)
	if err != nil {
		log.Fatalf("Could not marshall Google secrets: %s", err)
	}
	opt := option.WithCredentialsJSON(jsonCreds)

	return firebase.NewApp(context.Background(), nil, opt)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
