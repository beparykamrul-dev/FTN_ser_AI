package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	backend "github.com/beparykamrul-dev/FTN_ser_AI/backend"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	addr := os.Getenv("FTN_ADDR")
	if addr == "" { addr = "127.0.0.1:8080" }
	if err := backend.RunHTTP(ctx, addr); err != nil { log.Fatal(err) }
}
