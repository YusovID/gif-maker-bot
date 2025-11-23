package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/YusovID/gif-maker-bot/internal/config"
	"github.com/YusovID/gif-maker-bot/internal/converters"
	"github.com/YusovID/gif-maker-bot/internal/service"
	"github.com/YusovID/gif-maker-bot/internal/telegram"
	"github.com/YusovID/gif-maker-bot/pkg/logger/sl"
	"github.com/YusovID/gif-maker-bot/pkg/logger/slogadapter"
	"github.com/YusovID/gif-maker-bot/pkg/logger/slogpretty"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log := slogpretty.SetupLogger(cfg.Env)
	log.Info("starting gif-maker-bot", slog.String("env", cfg.Env))

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Error("BOT_TOKEN is not set")
		os.Exit(1)
	}

	// use adapter because tgbotapi use default logger
	err = tgbotapi.SetLogger(&slogadapter.SlogAdapter{Slog: log})
	if err != nil {
		log.Error("SetLogger failed", sl.Err(err))
		os.Exit(1)
	}

	b, err := telegram.NewBot(ctx, token, log)
	if err != nil {
		log.Error("bot initialization failed", sl.Err(err))
		os.Exit(1)
	}

	fs := telegram.NewTelegramFileStorage(b.Bot)
	converter := converters.NewFFMPEGConverter(log)
	videoService := service.NewVideoService(converter, fs)

	b.SetVideoService(videoService)

	wg := &sync.WaitGroup{}
	wg.Go(func() {
		err := b.Run()
		if err != nil {
			log.Error("bot starting failed", sl.Err(err))
			return
		}
	})

	<-ctx.Done()

	wg.Wait()

	log.Info("gif-maker-bot was stopped")
}
