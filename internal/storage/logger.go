package storage

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type LoggerStorage struct {
	logger *slog.Logger
}

func NewLoggerStorage(project_state string) *LoggerStorage {
	if project_state == "prod" {
		return &LoggerStorage{
			logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		}
	} else {
		return &LoggerStorage{
			logger: slog.New(tint.NewHandler(os.Stdout, &tint.Options{
				Level: slog.LevelDebug,
			})),
		}
	}

}

func (l *LoggerStorage) SetLogger() {
	slog.SetDefault(l.logger)
}
