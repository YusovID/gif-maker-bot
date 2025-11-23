package telegram

import (
	"github.com/YusovID/gif-maker-bot/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func toDomainMessage(tgMsg *tgbotapi.Message) (*domain.Message, bool) {
	if tgMsg.Video == nil {
		return nil, false
	}

	return &domain.Message{
		From: domain.User{
			ID:       tgMsg.From.ID,
			Username: tgMsg.From.UserName,
		},
		Chat: domain.Chat{
			ChatID: tgMsg.Chat.ID,
		},
		File: domain.File{
			ID:       tgMsg.Video.FileID,
			Type:     domain.Video,
			FileName: tgMsg.Video.FileName,
		},
		ReplyTo: tgMsg.MessageID,
	}, true
}
