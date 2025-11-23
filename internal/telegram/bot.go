package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/YusovID/gif-maker-bot/internal/apperrors"
	"github.com/YusovID/gif-maker-bot/internal/domain"
	"github.com/YusovID/gif-maker-bot/pkg/logger/sl"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type VideoService interface {
	VideoToGIF(ctx context.Context, msg *domain.Message) (string, error)
}

type Bot struct {
	Bot          *tgbotapi.BotAPI
	ctx          context.Context
	log          *slog.Logger
	videoService VideoService
	sem          chan struct{}
}

func NewBot(ctx context.Context, token string, log *slog.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("new bot init failed: %v", err)
	}

	bot.Debug = true

	log.Info("Authorized on account", slog.String("username", bot.Self.UserName))

	numWorkers, err := getWorkersNum()
	if err != nil {
		return nil, fmt.Errorf("getWorkersNum failed: %v", err)
	}

	return &Bot{
		Bot: bot,
		ctx: ctx,
		log: log,
		sem: make(chan struct{}, numWorkers),
	}, nil
}

func (b *Bot) SetVideoService(vs VideoService) {
	b.videoService = vs
}

func (b *Bot) Run() error {
	const op = "internal.bot.Run"
	log := b.log.With("op", op)
	log.Info("starting bot")

	err := b.notifyAdmin()
	if err != nil {
		return fmt.Errorf("notifyAdmin failed: %v", err)
	}

	b.Process()

	return nil
}

func (b *Bot) Process() {
	const op = "internal.bot.Process"
	log := b.log.With("op", op)
	log.Info("starting processing updates")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Bot.GetUpdatesChan(u)

	for {
		select {
		case <-b.ctx.Done():
			b.Bot.StopReceivingUpdates()
			return

		case update := <-updates:
			if update.Message == nil {
				continue
			}

			domainMsg, ok := toDomainMessage(update.Message)
			if !ok {
				err := b.replyWithError(b.ctx, update.Message.Chat.ID, update.Message.MessageID, apperrors.ErrUnsupportedFileType)
				if err != nil {
					log.Error("replyWithError failed", sl.Err(err))
				}

				continue
			}

			go func() {
				b.sem<-struct{}{}

				defer func() { <-b.sem }()

				err := b.ProcessVideo(domainMsg)
				if err != nil {
					log.Error("ProcessVideo failed", sl.Err(err))
					return
				}
			}()
		}
	}
}

func (b *Bot) ProcessVideo(msg *domain.Message) error {
	gifPath, err := b.videoService.VideoToGIF(b.ctx, msg)
	if err != nil {
		replyErr := b.replyWithError(b.ctx, msg.Chat.ChatID, msg.ReplyTo, err)
		if replyErr != nil {
			return fmt.Errorf("replyWithError failed: %w", err)
		}

		return fmt.Errorf("VideoToGIF failed: %w", err)
	}

	err = b.replyWithGIF(b.ctx, gifPath, msg.Chat.ChatID, msg.ReplyTo)
	if err != nil {
		replyErr := b.replyWithError(b.ctx, msg.Chat.ChatID, msg.ReplyTo, err)
		if replyErr != nil {
			return fmt.Errorf("replyWithError failed: %w", err)
		}

		return fmt.Errorf("replyWithGIF failed: %w", err)
	}

	return nil
}

func (b *Bot) notifyAdmin() error {
	adminIDStr := os.Getenv("ADMIN_TG_ID")
	if adminIDStr == "" {
		return fmt.Errorf("")
	}

	adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("admin id converting failed: %v", err)
	}

	msg := tgbotapi.NewMessage(adminID, "bot started")
	if _, err = b.Bot.Send(msg); err != nil {
		return fmt.Errorf("message sending failed: %v", err)
	}

	return nil
}

func (b *Bot) replyWithGIF(ctx context.Context, gifPath string, chatID int64, msgID int) error {
	msg := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(gifPath))
	msg.ReplyToMessageID = msgID

	_, err := b.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("gif sending failed: %v", err)
	}

	return nil
}

func (b *Bot) replyWithError(ctx context.Context, chatID int64, msgID int, intErr error) error {
	msg := tgbotapi.NewMessage(chatID, intErr.Error())
	msg.ReplyToMessageID = msgID

	_, err := b.Bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error message sending failed: %v", err)
	}

	return nil
}

func getWorkersNum() (int, error) {
	numWorkersStr := os.Getenv("WORKER_POOL_SIZE")
	if numWorkersStr == "" {
		return -1, fmt.Errorf("WORKER_POOL_SIZE not set")
	}

	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		return -1, fmt.Errorf("can't parse WORKER_POOL_SIZE: %v", err)
	}

	if numWorkers <= 0 {
		return -1, fmt.Errorf("WORKER_POOL_SIZE can't be less than or equal to zero")
	}

	return numWorkers, nil
}
