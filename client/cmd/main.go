package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "github.com/horirinn2000/grpc-connect-envoy/proto/gen/greet/v1"
	"github.com/horirinn2000/grpc-connect-envoy/proto/gen/greet/v1/greetv1connect"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://envoy:8080",
	)
	res, err := client.Greet(
		context.Background(),
		&greetv1.GreetRequest{Name: "Jane"},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Greeting)
}
