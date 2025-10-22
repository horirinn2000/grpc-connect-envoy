package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/horirinn2000/grpc-connect-envoy/services/auth/internal/auth"
	"github.com/horirinn2000/grpc-connect-envoy/services/auth/internal/config"
	"github.com/horirinn2000/grpc-connect-envoy/services/auth/internal/key"
	authhandler "github.com/horirinn2000/grpc-connect-envoy/services/auth/pkg/auth_handler"
)

type requestTracker struct {
	wg *sync.WaitGroup
	h  http.Handler
}

func (t *requestTracker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.wg.Add(1)
	defer t.wg.Done()
	t.h.ServeHTTP(w, r)
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	privateKey, err := key.LoadPrivateKey(cfg.JWT.PrivateKeyPath)
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}

	svc := auth.NewService(privateKey, cfg.Auth.Username, cfg.Auth.Password, cfg.JWT.TokenExp)
	path, handler := authhandler.NewHandler(svc)

	var wg sync.WaitGroup
	trackedHandler := &requestTracker{wg: &wg, h: handler}

	mux := http.NewServeMux()
	mux.Handle(path, trackedHandler)

	p := new(http.Protocols)
	p.SetHTTP1(true)
	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	s := http.Server{
		Addr:      "0.0.0.0:" + cfg.Server.Port,
		Handler:   mux,
		Protocols: p,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("GRPC Server Start\n")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v\n", err)
		}
	}()

	//graceful shutdown
	log.Printf("SIGNAL SIGTERM:%d received, then shutting down...\n", <-signalChan)

	shutdownTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("all in-flight requests finished")
	case <-ctx.Done():
		log.Println("graceful timeout reached; some requests may be terminated")
	}

	log.Println("shutdown complete")
}
