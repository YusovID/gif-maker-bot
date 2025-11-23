package converters

import (
	"strconv"

	"github.com/YusovID/gif-maker-bot/internal/domain"
)

type FileData struct {
	UserID   string
	FileID   string
	FileName string
}

func NewFileData(msg *domain.Message) *FileData {
	userID := strconv.Itoa(int(msg.From.ID))

	return &FileData{
		UserID:   userID,
		FileID:   msg.File.ID,
		FileName: msg.File.FileName,
	}
}
