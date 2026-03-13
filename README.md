# logentrycheck

Go-линтер для проверки лог-записей, совместимый с `golangci-lint`.

Линтер проверяет вызовы логгеров `log/slog` и `go.uber.org/zap` на соответствие установленным правилам.

## Содержание

- [Проверяемые правила](#проверяемые-правила)
- [Установка](#установка)
- [Запуск как standalone-инструмент](#запуск-как-standalone-инструмент)
- [Интеграция с golangci-lint](#интеграция-с-golangci-lint)
- [Конфигурация](#конфигурация)
- [Авто-исправление](#авто-исправление)
- [Примеры использования](#примеры-использования)

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

## Установка

### Требования

- Go 1.22+
- git

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

### Конфигурация через golangci-lint

При использовании плагина параметры передаются через `.golangci.yml`:
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
        settings:
          # отключить проверку строчной буквы
          disableLowercase: true
          # отключить проверку английского языка
          disableEnglish: true
          # отключить проверку спецсимволов
          disableSpecialChars: true
          # отключить проверку чувствительных данных
          disableSensitive: true
          # кастомные ключевые слова (заменяют список по умолчанию)
          keywords:
            - "mytoken"
            - "internalkey"
            - "bearer"
```

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

## Примеры использования

В качестве проектов, для которого проводилось тестирование применения линтера был выбран https://github.com/uber-go/zap, так как в нём присутствуют проверяемые вызовы.

### Standalone запуск

На начальном этапе был установлен и запущен линтер в standalone режиме. Результат выполнения представлен ниже:
![Результат standalone запуска logentrycheck](<.github/images/Screenshot 2026-03-13 at 17.24.30.png>)
Изображение показывает, что линтер успешно детектировал все замечания.

### Запуск плагина для `golangci-lint`

Затем в корень директории были добавлены файлы `.custom-gcl.yml` и `.golangci.yml` с базовыми настройками сборки и конфигурации плагина. Плагин был успешно собран и запущен на ранее упомянутом проекте. Результат вывода:
![Результат запуска плагина для golangci-lint на проекте uber-go/zap](<.github/images/Screenshot 2026-03-13 at 17.28.39.png>)
Вывод на изображении соответствует ранее полученному.

### Конфигурация плагина

Далее в файле `.golangci.yml` была отключена проверка одного из правил. Листинг данного файла приведен ниже.
![Конфигурация .golangci.yml с отключённым правилом](<.github/images/Screenshot 2026-03-13 at 17.44.08.png>)
Плагин был повторно собран и запущен линтинг. Полученный результат не содержит данных о замечаниях, касающихся проверки исключенного правила.
![Результат запуска с отключённым правилом](<.github/images/Screenshot 2026-03-13 at 17.44.33.png>)

### Проверка автоисправления

Для исправления детектированных замечаний был запущен линтинг с флагом `--fix`. Результат приведен на рисунке снизу.
![Результат запуска с флагом --fix](<.github/images/Screenshot 2026-03-13 at 17.58.07.png>)

Вывод соответствует ожидаемому и подтверждается измененным кодом на указанных строках.
![Изменённый код после автоисправления](<.github/images/Screenshot 2026-03-13 at 17.58.27.png>)


### Проверка конфигурации ключевых слов

В силу того, что используемый проект не содержит логов, сообщения которых содержат "чувствительные" данные, противоречащие правилу, используемый плагин был использован для личного пет-проекта, доступного по адресу: https://github.com/IlyaChern12/rtce.

В корне проекта был создан файл следующего содержания, заведомо содержащий нарушения указанного правила:
![Тестовый файл с намеренными нарушениями правила sensitive](<.github/images/Screenshot 2026-03-13 at 17.47.24.png>)

После запуска линтера замечания были успешно детектированы:
![Детектированные нарушения правила sensitive с дефолтными паттернами](<.github/images/Screenshot 2026-03-13 at 17.54.04.png>)

Затем конфигурация линтера была изменена на обнаружение "чувствительных" данных по другому паттерну:
![Конфигурация кастомных ключевых слов в .golangci.yml](<.github/images/Screenshot 2026-03-13 at 17.55.06.png>)

Плагин был пересобран и повторно запущен:
![Результат запуска с кастомными ключевыми словами](<.github/images/Screenshot 2026-03-13 at 17.55.34.png>)

### Вывод

Таким образом, полученные результаты свидетельствуют об успешном применении разработанного линтера как в качестве standalone-инструмента, так и в качестве плагина для `golangci-lint` на реальных проектах. Кроме того, успешно реализована возможность конфигурации проверяемых правил, используемых паттернов обнаружения "чувствительных данных" и автоисправление обнаруженных замечаний.