package main

import (
	"log"
	"pratbacknd/internal/server"
)

func main() {
	server, err := server.New(server.Config{Port: "5000"})
	if err != nil {
		log.Fatalf("Could not create server : %s", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("Could not start the server : %s", err)
	}
}
