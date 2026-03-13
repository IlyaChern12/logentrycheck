package english

import (
	"context"
	"fmt"
	"log/slog"
)

type dummy struct{}

func (d dummy) Info(msg string) {}

func examples() {
	ctx := context.Background()

	slog.Info("user logged in")                     // ok
	slog.Info("пользователь вошёл")                 // want `log message must be in English only`
	slog.Info("user logged in 🎉")                   // want `log message must be in English only`
	slog.Info("üser alles")                         // want `log message must be in English only`
	slog.Debug("request started")                   // ok
	slog.InfoContext(ctx, "ошибка подключения")     // want `log message must be in English only`
	slog.InfoContext(ctx, "connection established") // ok

	// msg == ""
	slog.Info("")

	// found == false
	fmt.Println("hello")

	// found == false
	var d dummy
	d.Info("test")

	log := slog.Default()
	log.Info("test message")

	// pointer logger
	ptr := slog.Default()
	ptr.Info("another message")

	// context logger
	log.InfoContext(ctx, "context message")

	slog.Info(fmt.Sprintf("user %d logged in", 1))
}
