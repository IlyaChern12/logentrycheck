package english

import (
	"context"
	"log/slog"
)

func examples() {
	ctx := context.Background()

	slog.Info("user logged in")                     // ок
	slog.Info("пользователь вошёл")                 // want `log message must be in English only: "пользователь вошёл"`
	slog.Info("user logged in 🎉")                   // want `log message must be in English only: "user logged in 🎉"`
	slog.Info("über alles")                         // want `log message must be in English only: "über alles"`
	slog.Debug("request started")                   // ок
	slog.InfoContext(ctx, "ошибка подключения")     // want `log message must be in English only: "ошибка подключения"`
	slog.InfoContext(ctx, "connection established") // ок
}
