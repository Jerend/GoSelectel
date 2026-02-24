package lowercase

import "log/slog"

func TestLowercaseStart() {
	slog.Info("Starting server on port 8080")   // want "The message in the log starts with a capital letter"
	slog.Error("Failed to connect to database") // want "The message in the log starts with a capital letter"

	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
}
