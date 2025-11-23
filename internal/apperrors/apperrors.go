package apperrors

import "errors"

var (
	ErrUnsupportedFileType = errors.New("this file type does not support")
	ErrFileNotFound        = errors.New("file not found")
)
