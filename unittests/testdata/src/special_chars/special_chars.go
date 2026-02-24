package special_chars

import "log/slog"

func TestNoSpecialChars() {
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")

	slog.Info("server started!🚀")                 // want "log message contains invalid character"
	slog.Error("connection failed!!!")            // want "log message contains invalid character"
	slog.Warn("warning: something went wrong...") // want "log message contains invalid character"
}
