package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rayikume/payment-splitter/config"
	"github.com/rayikume/payment-splitter/internal/handlers"
	"github.com/rayikume/payment-splitter/internal/middleware"
	"github.com/rayikume/payment-splitter/internal/services"
)

func main() {
	cnfg := config.Load()
	port := cnfg.AppPort
	fmt.Println(port)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	splitService := services.NewSplitService()
	splitHandler := handlers.NewSplitHandler(splitService)

	mux := http.NewServeMux()
	splitHandler.RegisterRoutes(mux)

	var p http.Handler = mux
	p = middleware.Logger(p)
	p = middleware.RequestID(p)
	p = middleware.CORS(p)
	p = middleware.Recovery(p)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      p,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "port", port)
		errCh <- server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)

	select {
	case sig := <-quit:
		slog.Info("shutdown signal received", "signal", sig)
	case err := <-errCh:
		slog.Error("server error", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped gracefully")
}
