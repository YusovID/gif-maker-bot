package slogadapter

import (
	"fmt"
	"log/slog"
	"strings"
)

type SlogAdapter struct {
	Slog *slog.Logger
}

func (a *SlogAdapter) Println(v ...interface{}) {
	a.Slog.Debug(strings.TrimSpace(fmt.Sprint(v...)))
}

func (a *SlogAdapter) Printf(format string, v ...interface{}) {
	a.Slog.Debug(strings.TrimSpace(fmt.Sprintf(format, v...)))
}
