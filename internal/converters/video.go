package converters

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FFMPEGConverter struct {
	log *slog.Logger
}

func NewFFMPEGConverter(log *slog.Logger) *FFMPEGConverter {
	return &FFMPEGConverter{log: log}
}

func (c *FFMPEGConverter) VideoToGIF(ctx context.Context, fileStream io.ReadCloser, fileData *FileData) (string, error) {
	const op = "internal.bot.VideoToGIF"

	log := c.log.With("op", op)

	log.Info("starting video to gif convertation")

	videoPath := getVideoPath(fileData)
	videoDir := filepath.Dir(videoPath)

	gifPath := getGIFPath()
	gifDir := filepath.Dir(gifPath)

	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return "", fmt.Errorf("video directory creation failed: %v", err)
	}
	defer os.RemoveAll(videoDir)

	if err := os.MkdirAll(gifDir, 0755); err != nil {
		return "", fmt.Errorf("video directory creation failed: %v", err)
	}

	tempFile, err := os.Create(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, fileStream)
	if err != nil {
		return "", fmt.Errorf("failed to write stream to file: %v", err)
	}

	var errorBuf bytes.Buffer

	err = ffmpeg.Input(videoPath).Output(gifPath).OverWriteOutput().
		WithErrorOutput(&errorBuf).Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w | Log: %s", err, errorBuf.String())
	}

	log.Info("gif saved", slog.String("path", gifPath))

	return gifPath, nil
}

func getVideoPath(fileData *FileData) string {
	fileExt := filepath.Ext(fileData.FileName)
	newFileName := uuid.New().String() + fileExt

	return filepath.Join(os.TempDir(), GIFMakerWorkspaceDir, fileData.UserID, newFileName)
}

func getGIFPath() string {
	newFileName := uuid.New().String() + ".gif"

	return filepath.Join(os.TempDir(), GIFMakerWorkspaceDir, ResultDir, newFileName)
}
