package lowercase

import (
	"context"
	"log/slog"
)

func examples() {
	ctx := context.Background()

	slog.Info("user logged in")           // ок
	slog.Info("User logged in")           // want `log message should start with a lowercase letter: "User logged in"`
	slog.Info("REQUEST failed")           // want `log message should start with a lowercase letter: "REQUEST failed"`
	slog.Debug("debug message")           // ок
	slog.Debug("Debug message")           // want `log message should start with a lowercase letter: "Debug message"`
	slog.Warn("something went wrong")     // ок
	slog.Error("Error occurred")          // want `log message should start with a lowercase letter: "Error occurred"`
	slog.InfoContext(ctx, "user updated") // ок
	slog.InfoContext(ctx, "User updated") // want `log message should start with a lowercase letter: "User updated"`
}
