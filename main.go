package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/rayikume/payment-splitter/config"
)

func main() {
	cnfg := config.Load()
	port := cnfg.AppPort
	fmt.Println(port)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}
