run:
	echo "Run triggered"
	go run main.go

build:
	echo "Building for linux"
	env GOOS=linux GOARCH=amd64 go build -o bin/api api/main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/hello testLambda/main.go

deploy: build
	serverless deploy --aws-profile adil