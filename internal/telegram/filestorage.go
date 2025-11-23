package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramFileStorage struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramFileStorage(bot *tgbotapi.BotAPI) *TelegramFileStorage {
	return &TelegramFileStorage{bot: bot}
}

func (tfs *TelegramFileStorage) Get(ctx context.Context, fileID string) (io.ReadCloser, error) {
	fileURL, err := tfs.bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, fmt.Errorf("can't get file url: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}

	resp, err := tfs.bot.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %v", err)
	}

	return resp.Body, nil
}
