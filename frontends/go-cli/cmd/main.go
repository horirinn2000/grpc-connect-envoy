package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	authv1 "github.com/horirinn2000/grpc-connect-envoy/services/auth/api/gen/auth/v1"
	"github.com/horirinn2000/grpc-connect-envoy/services/auth/api/gen/auth/v1/authv1connect"
	greetv1 "github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1"
	"github.com/horirinn2000/grpc-connect-envoy/services/greet/api/gen/greet/v1/greetv1connect"
	thanksv1 "github.com/horirinn2000/grpc-connect-envoy/services/thanks/api/gen/thanks/v1"
	"github.com/horirinn2000/grpc-connect-envoy/services/thanks/api/gen/thanks/v1/thanksv1connect"
)

type authInterceptor struct {
	token string
}

func (i *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		req.Header().Set("Authorization", "Bearer "+i.token)
		return next(ctx, req)
	}
}

func (i *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

func main() {
	// 1. Authenticate and get a token
	authClient := authv1connect.NewAuthServiceClient(
		http.DefaultClient,
		"http://envoy-client:8080",
	)
	authRes, err := authClient.Authenticate(
		context.Background(),
		&authv1.AuthenticateRequest{
			Username: "user",
			Password: "password",
		},
	)
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}
	log.Println("Successfully authenticated")

	// 2. Call Greet service with the token
	interceptor := &authInterceptor{token: authRes.Token}
	greetClient := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://envoy-client:8080",
		connect.WithInterceptors(interceptor),
	)
	resGreet, err := greetClient.Greet(
		context.Background(),
		&greetv1.GreetRequest{
			Name: "horirinn",
		},
	)
	if err != nil {
		log.Println("Greet error:", err)
	} else {
		log.Println("Greet response:", resGreet.Greeting)
	}

	// 3. Call Thanks service without a token
	thanksClient := thanksv1connect.NewThanksServiceClient(
		http.DefaultClient,
		"http://envoy-client:8080",
	)
	resThanks, err := thanksClient.Thanks(
		context.Background(),
		&thanksv1.ThanksRequest{
			Name: "horirinn",
		},
	)
	if err != nil {
		log.Println("Thanks error:", err)
	} else {
		log.Println("Thanks response:", resThanks.Thanks)
	}

	// 4. Refresh
	// 1. Refresh
	_, err = authClient.RefreshToken(
		context.Background(),
		&authv1.RefreshTokenRequest{
			RefreshToken: authRes.RefreshToken,
		},
	)
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}
	log.Println("Successfully refresh")
}
