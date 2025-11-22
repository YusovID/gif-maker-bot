package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/YusovID/gif-maker-bot/internal/bot"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalf("BOT_TOKEN is not set\n")
	}

	b, err := bot.NewBot(ctx, token)
	if err != nil {
		log.Fatalf("bot initialization failed: %v\n", err)
	}

	wg := &sync.WaitGroup{}
	wg.Go(func() { b.Run() })

	<-ctx.Done()

	wg.Wait()

	log.Println("stopping bot")
}
