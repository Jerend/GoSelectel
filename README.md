# GoSelectel
### mylinter - анализатор лог-сообщений

Проверяет лог-сообщения в `log/slog` и `go.uber.org/zap` на соответствие правилам:
- Лог-сообщения должны начинаться со строчной буквы
- Лог-сообщения должны быть только на английском языке
- Лог-сообщения не должны содержать спецсимволы или эмодзи
- Лог-сообщения не должны содержать потенциально чувствительные данные

## Сборка и запуск

### Требования

- **Go** версии 1.26.0
- **golangci-lint** версии 2.10.1
- **Используемая среда:** WSL, macOS или Linux

### 1. Клонируйте репозиторий и перейдите в него

```bash
git clone https://github.com/Jerend/GoSelectel.git
cd GoSelectel
```
### 2. Соберите плагин
```bash
CGO_ENABLED=1 go build -buildmode=plugin -o ./plugin/mylinter.so ./plugin/plugin.go
```

### 3. Перейдите в корневую директорию проекта, который хотите проверить, и создайте файл конфигурации
```bash
cd /путь/к/проекту
touch .golangci.yml
```

### 4. Откройте файл .golangci.yml и добавьте следующую конфигурацию. В поле path укажите путь к собранному файлу плагина
```bash
version: "2"
linters:
  enable:
    - mylinter
  settings:
    custom:
      mylinter:
        type: goplugin
        path: /путь/к/GoSelectel/plugin/mylinter.so
```

### 5. Запустите анализ

### Запуск линтера без исправления ошибок
```bash
golangci-lint run ./...
```
### Запуск линтера с исправлением ошибок
```bash
golangci-lint run --fix ./...
```

Исправления при запуске с флагом `--fix`:
- Заглавная буква в начале → замена на строчную
- Недопустимые символы → удаление
- Не-английские буквы → удаление
- Чувствительные данные → замена на [change]

## Примеры использования
### При запуске линтера на следующих командах:
```bash
slog.Info("user password: " + password)
slog.Debug("apiKey: " + apiKey)
slog.Info("token: " + token)
```
### Вывод без --fix:
```bash
potential sensitive data "password" is being logged directly
potential sensitive data "apiKey" is being logged directly
potential sensitive data "token" is being logged directly
```
### Вывод с помощью --fix:
```bash
slog.Info("user password: " + "[change]")
slog.Debug("apiKey: " + "[change]")
slog.Info("token: " + "[change]")
```

### При запуске линтера на следующих командах:
```bash
slog.Info("Starting server on port 8080")
slog.Error("Failed to connect to database")
```
### Вывод без --fix:
```bash
The message in the log starts with a capital letter, got "Starting server on port 8080"
The message in the log starts with a capital letter, got "Failed to connect to database"
```
### Вывод с помощью --fix:
```bash
slog.Info("starting server on port 8080")
slog.Error("failed to connect to database")
```