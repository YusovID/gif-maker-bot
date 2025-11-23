package service

import (
	"context"
	"fmt"
	"io"

	converter "github.com/YusovID/gif-maker-bot/internal/converters"
	"github.com/YusovID/gif-maker-bot/internal/domain"
)

type Converter interface {
	VideoToGIF(ctx context.Context, fileStream io.ReadCloser, fileData *converter.FileData) (string, error)
}

type VideoService struct {
	converter Converter
	fs        FileStorage
}

func NewVideoService(converter Converter, fileStorage FileStorage) *VideoService {
	return &VideoService{
		converter: converter,
		fs:        fileStorage,
	}
}

func (vs *VideoService) VideoToGIF(ctx context.Context, msg *domain.Message) (string, error) {
	fileStream, err := vs.fs.Get(ctx, msg.File.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get file stream: %w", err)
	}
	defer fileStream.Close()

	fileData := converter.NewFileData(msg)

	gifPath, err := vs.converter.VideoToGIF(ctx, fileStream, fileData)
	if err != nil {
		return "", fmt.Errorf("service.VideoToGIF failed: %v", err)
	}

	return gifPath, nil
}
