package sensitive

import "log/slog"

func TestNoSensitiveVars(password, apiKey, token string) {
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")

	slog.Info("user password: " + password) // want "potential sensitive data \"password\" is being logged directly"

	slog.Debug("apiKey: " + apiKey) // want "potential sensitive data \"apiKey\" is being logged directly"

	slog.Info("token: " + token) // want "potential sensitive data \"token\" is being logged directly"
}
