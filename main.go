package main

import (
	"fmt"

	"github.com/rayikume/payment-splitter/config"
)

type Config struct {
	AppPort string
}

func main() {
	cnfg := config.Load()
	port := cnfg.AppPort
	fmt.Println(port)
}
