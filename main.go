package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/rayikume/payment-splitter/config"
	"github.com/rayikume/payment-splitter/internal/handlers"
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
}
