package english_only

import "log/slog"

func TestEnglishOnly() {
	slog.Info("starting server")
	slog.Error("failed to connect to database")

	slog.Info("запуск сервера")                    // want "log message contains non-English characters"
	slog.Error("ошибка подключения к базе данных") // want "log message contains non-English characters"
}
