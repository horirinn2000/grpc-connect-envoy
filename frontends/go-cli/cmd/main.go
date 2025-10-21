package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1"
	"github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1/greetv1connect"
	thanksv1 "github.com/horirinn2000/grpc-connect-envoy/services/thanks/api/gen/thanks/v1"
	"github.com/horirinn2000/grpc-connect-envoy/services/thanks/api/gen/thanks/v1/thanksv1connect"
)

func main() {
	callGreet()
	callThanks()
}

func callGreet() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://envoy-client:8080",
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

func callThanks() {
	client := thanksv1connect.NewThanksServiceClient(
		http.DefaultClient,
		"http://envoy-client:8080",
	)
	res, err := client.Thanks(
		context.Background(),
		&thanksv1.ThanksRequest{Name: "Jane"},
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Thanks)
}
