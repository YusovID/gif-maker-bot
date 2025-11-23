package service

import (
	"context"
	"io"
)

type FileStorage interface {
	Get(ctx context.Context, fileID string) (io.ReadCloser, error)
}
