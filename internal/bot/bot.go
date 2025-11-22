package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Bot *tgbotapi.BotAPI
	ctx context.Context
}

func NewBot(ctx context.Context, token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("new bot init failed: %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	return &Bot{
		Bot: bot,
		ctx: ctx,
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Bot.GetUpdatesChan(u)

	for {
		select {
		case <-b.ctx.Done():
			b.Bot.StopReceivingUpdates()
			return

		case update := <-updates:
			if update.Message != nil {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				if _, err := b.Bot.Send(msg); err != nil {
					log.Printf("ERROR: message sending failed: %v\n", err)
				}
			}
		}
	}
}
