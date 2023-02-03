run:
	echo "Run triggered"
	go run main.go

build:
	echo "Building for linux"
	env GOOS=linux GOARCH=amd64 go build -o bin/api api/main.go

deploy: build
	serverless deploy --param="allowedOrigin=https://master.d14f8mlnk4lkw2.amplifyapp.com/" --aws-profile adil

deploy_dev: build
	serverless deploy --param="allowedOrigin=http://localhost:5173" --aws-profile adil --param="stage=dev"

genmocks:
	mockgen -source=internal/storage/storage.go -destination=internal/storage/mock.go -package=storage
	mockgen -source=internal/utils/uuidgenerator.go -destination=internal/utils/uuidgeneratormock.go -package=utils

#--param="allowedOrigin=https://master.d14f8mlnk4lkw2.amplifyapp.com/"