run:
	echo "Run triggered"
	go run main.go

build:
	echo "Building for linux"
	env GOOS=linux GOARCH=amd64 go build -o bin/api api/main.go

deploy: build
	serverless deploy --aws-profile adil --force

#--param="allowedOrigin=http://localhost:5173/"
#--param="allowedOrigin=https://master.d14f8mlnk4lkw2.amplifyapp.com/"