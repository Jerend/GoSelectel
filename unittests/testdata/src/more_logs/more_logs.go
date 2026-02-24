package more_logs

import (
	"log/slog"

	"go.uber.org/zap"
)

func Alltests(password string, apiKey string, token string, secret string, auth_token string) {
	zapLogger, _ := zap.NewProduction()
	sugar := zapLogger.Sugar()

	slog.Info("starting server")
	slog.Error("failed to connect to database")
	slog.Info("запуск сервера")
	slog.Error("Ошибка подключения!!! к базе данных")
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")
	slog.Debug("api_key=" + apiKey)
	slog.Info("token: " + token)
	slog.Info("server started!🚀")
	slog.Error("connection failed!!!")
	slog.Warn("warning: something went wrong...")
	zapLogger.Info("Token is valid! ✔️")
	zapLogger.Debug("Processing request!!!")
	zapLogger.Debug("Processing request!!!")
	sugar.Info("Processing payment")
	sugar.Info("Backup size: 1.2GB...")
	zapLogger.Warn("Payment validation")
	slog.Info("Миграция базы данных выполнена успешно")
	slog.Warn("connection start")
	zapLogger.Error("недостаточно места на диске")
	zapLogger.Info("Ротация логов выполнена: архив logs-20250223.tar.gz создан")
	zapLogger.Info("secret: " + secret)
	slog.Info("auth_token: " + auth_token)
}
