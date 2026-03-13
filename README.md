# logentrycheck

Go-линтер для проверки лог-записей, совместимый с `golangci-lint`.

Линтер проверяет вызовы логгеров `log/slog` и `go.uber.org/zap` на соответствие установленным правилам.

## Содержание

- [Установка](#установка)
- [Запуск как standalone-инструмент](#запуск-как-standalone-инструмент)
- [Интеграция с golangci-lint](#интеграция-с-golangci-lint)
- [Конфигурация](#конфигурация)
- [Авто-исправление](#авто-исправление)
- [Примеры использования](#примеры-использования)
- [Правила](#правила)

## Установка

### Требования

- Go 1.22+

### Сборка бинарника

```bash
git clone https://github.com/IlyaChern12/logentrycheck
cd logentrycheck
go build -o logentrycheck ./cmd/logentrycheck
```

## Обычный запуск (standalone)

```bash
# установка
go install github.com/IlyaChern12/logentrycheck/cmd/logentrycheck@latest

# запуск
logentrycheck ./...
```

## Интеграция с golangci-lint

Линтер поддерживает [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/) для `golangci-lint`.

### Шаг 1: Создать `.custom-gcl.yml` в корне проекта

```yaml
version: v2.11.3
plugins:
  # подключение через go proxy
  - module: 'github.com/IlyaChern12/logentrycheck'
    version: v0.1.0

  # или из локального источника
  - module: 'github.com/IlyaChern12/logentrycheck'
    path: <path>/logentrycheck
```

### Шаг 2: Создать `.golangci.yml`

```yaml
version: "2"

linters:
  default: none
  enable:
    - logentrycheck
  settings:
    custom:
      logentrycheck:
        type: "module"
        description: checks log entries for common mistakes in log messages
```

### Шаг 3: Собрать кастомный бинарник

```bash
golangci-lint custom -v
```

### Шаг 4: Запустить

```bash
./custom-gcl run ./...
```

### Интеграции

#### GitHub Actions

```yaml
- name: Build custom golangci-lint
  run: golangci-lint custom -v

- name: Run logentrycheck
  run: ./custom-gcl run ./...
```

#### GitLab CI
```yaml
lint:
  stage: lint
  image: golang:1.24
  before_script:
    - curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.3
  script:
    - golangci-lint custom -v
    - ./custom-gcl run ./...
```

---

## Конфигурация

### Отключение отдельных правил

Правила можно отключать через флаги:

```bash
# отключить проверку строчной буквы
logentrycheck -logentrycheck.disable-lowercase ./...

# отключить проверку английского языка
logentrycheck -logentrycheck.disable-english ./...

# отключить проверку спецсимволов
logentrycheck -logentrycheck.disable-special-chars ./...

# отключить проверку чувствительных данных
logentrycheck -logentrycheck.disable-sensitive ./...
```

### Кастомные ключевые слова для правила sensitive

По умолчанию используется встроенный список ключевых слов. При указании кастомных слов они заменяют список по умолчанию:

```bash
logentrycheck -logentrycheck_sensitive.keywords="mytoken,internalkey,bearer" ./...
```

---

## Авто-исправление

Линтер поддерживает автоматическое исправление для правил:

- **lowercase** — первая буква сообщения приводится к нижнему регистру
- **special_chars** — запрещённые символы удаляются из сообщения

```bash
logentrycheck -fix ./...
```

Пример:

```go
// до
slog.Info("Starting server")
slog.Info("connection failed!!!")

// после
slog.Info("starting server")
slog.Info("connection failed")
```

## Проверяемые правила

### 1. Строчная буква в начале сообщения (`logentrycheck_lowercase`)

Лог-сообщения должны начинаться со строчной буквы.

```go
// неправильно
slog.Info("Starting server on port 8080")
zap.Error("Failed to connect to database")

// правильно
slog.Info("starting server on port 8080")
zap.Error("failed to connect to database")
```

### 2. Только английский язык (`logentrycheck_english`)

Лог-сообщения должны быть написаны только на английском языке.

```go
// неправильно
slog.Info("запуск сервера")
zap.Error("ошибка подключения к базе данных")

// правильно
slog.Info("starting server")
zap.Error("failed to connect to database")
```

### 3. Запрет спецсимволов и эмодзи (`logentrycheck_special_chars`)

Лог-сообщения не должны содержать спецсимволы или эмодзи. Разрешены: буквы, цифры, пробел и символы ` - _ / . , ( ) [ ] { } @ # % + = < > :`.

```go
// неправильно
slog.Info("server started! 🚀")
zap.Error("connection failed!!!")
slog.Warn("what happened?")

// правильно
slog.Info("server started")
zap.Error("connection failed")
slog.Warn("something went wrong")
```

### 4. Запрет чувствительных данных (`logentrycheck_sensitive`)

Лог-сообщения не должны содержать потенциально чувствительные данные через конкатенацию строк.

Ключевые слова по умолчанию: `password`, `passwd`, `secret`, `token`, `api_key`, `apikey`, `auth`, `credential`, `private_key`, `access_key`, `session`.

```go
// неправильно
slog.Info("user password: " + password)
zap.Debug("api_key=" + apiKey)
slog.Info("token: " + token)

// правильно
slog.Info("user authenticated successfully")
zap.Debug("api request completed")
slog.Info("token validated")
```