package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apiservices/youtube-thumbnail-api/internal/youtube/api"
	"apiservices/youtube-thumbnail-api/internal/youtube/auth"
	"apiservices/youtube-thumbnail-api/internal/youtube/service"
)

func main() {
	logger := log.New(os.Stdout, "[youtube-thumbnail] ", log.LstdFlags)

	port := envString("PORT", "8099")
	apiKey := envString("YOUTUBE_API_KEY", "dev-youtube-key")
	environment := envString("ENVIRONMENT", "development")
	proxySecret := envString("RAPIDAPI_PROXY_SECRET", "")

	if apiKey == "dev-youtube-key" {
		logger.Println("YOUTUBE_API_KEY not set, using default development key")
	}

	if environment == "production" && proxySecret == "" {
		logger.Println("WARNING: RAPIDAPI_PROXY_SECRET not set in production mode")
	}

	svc := service.NewService()
	handler := api.NewHandler(svc)

	mux := http.NewServeMux()

	// Wrap handler with auth middleware that includes environment context
	authHandler := auth.Middleware(apiKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add environment context to request headers for middleware
		if environment == "production" {
			r.Header.Set("X-Environment", "production")
			r.Header.Set("X-Expected-Proxy-Secret", proxySecret)
		}
		handler.ServeHTTP(w, r)
	}))

	mux.Handle("/v1/youtube/", authHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Printf("service listening on :%s (environment: %s)", port, environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("shutdown error: %v", err)
	}
}

func envString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
