package sensitive

import (
	"log/slog"
)

func examples() {
	password := "secret123"
	apiKey := "key123"
	token := "tok123"

	// plain literals
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")
	slog.Info("user logged in")
	slog.Debug("request completed")

	// concatenations with sensitive keyword
	slog.Info("user password: " + password) // want `log message may contain sensitive data \(keyword: "password"\): "user password: "`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data \(keyword: "api_key"\): "api_key="`
	slog.Info("token: " + token)            // want `log message may contain sensitive data \(keyword: "token"\): "token: "`
}
