package log

import (
	"context"
	"log/slog"
	"os"
	"time"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func Info(message string, params ...interface{}) {
	InfoContext(context.Background(), message, params...)
}

func InfoContext(ctx context.Context, message string, params ...interface{}) {
	logger.InfoContext(ctx, message, params...)
}

func Error(message string, err error, params ...interface{}) {
	ErrorContext(context.Background(), message, err, params...)
}

func ErrorContext(ctx context.Context, message string, err error, params ...interface{}) {
	record := slog.NewRecord(time.Now(), slog.LevelError, message, 0)
	record.Add("error", err)
	record.Add(params...)
	_ = logger.Handler().Handle(ctx, record)
}

func Warn(message string, params ...interface{}) {
	WarnContext(context.Background(), message, params...)
}

func WarnContext(ctx context.Context, message string, params ...interface{}) {
	logger.WarnContext(ctx, message, params...)
}
