# `cli` — Минималистичный CLI-фреймворк для Go-приложений

[![Go CI](https://github.com/shuldan/cli/workflows/Go%20CI/badge.svg)](https://github.com/shuldan/cli/actions)
[![codecov](https://codecov.io/gh/shuldan/cli/branch/main/graph/badge.svg)](https://codecov.io/gh/shuldan/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/shuldan/cli)](https://goreportcard.com/report/github.com/shuldan/cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Пакет `cli` предоставляет простой, но расширяемый каркас для создания консольных утилит на Go с декларативным описанием команд, автоматической генерацией справки, middleware-цепочками и безопасным выполнением.

---

## 🚀 Основные возможности

- **Stateless-команды** — аргументы и опции описываются декларативно, распарсенный ввод передаётся через `*Input`. Команды потокобезопасны по дизайну.
- **Разделение аргументов и опций** — позиционные аргументы (обязательные/необязательные) и именованные опции с типизированными значениями.
- **Автоматическая справка** — встроенные команды `help` и `version`. Поддержка `--help` / `-h` и `--version` / `-v` как глобальных флагов.
- **Группировка команд** — логические группы с детерминированным порядком вывода.
- **Middleware** — цепочки обработчиков для логирования, замера времени, авторизации и других сквозных задач.
- **Типизированные ошибки и exit-коды** — `CommandNotFoundError`, `MissingArgumentError`, `PanicError` и другие. Каждая ошибка несёт код завершения процесса.
- **Безопасное выполнение** — перехват паник с сохранением полного стектрейса.
- **Graceful shutdown** — встроенный хелпер `SignalContext` для обработки SIGINT/SIGTERM.
- **Тестируемость** — полная поддержка внедрения зависимостей через `io.Reader`, `io.Writer`, `context.Context`.

---

## 📦 Установка

```sh
go get github.com/shuldan/cli
```

Требуется Go 1.24+.

### Инструменты разработки

```sh
make install-tools
```

Устанавливает:
- `golangci-lint` (v2.4.0)
- `goimports`
- `gosec`

---

## 🛠️ Работа с проектом

| Команда | Описание |
|---------|----------|
| `make all` | Форматирование, линтер, security-сканирование, тесты |
| `make ci` | CI-пайплайн: форматирование, `go vet`, линтер, тесты с покрытием |
| `make fmt` | Форматирование кода и сортировка импортов |
| `make test` | Запуск тестов |
| `make test-coverage` | Тесты с отчётом о покрытии |

---

## 🧱 Архитектура

### Интерфейс `Command`

Каждая команда реализует единый интерфейс. Команда полностью stateless — вся конфигурация декларативна, а распарсенные данные передаются в `Execute` через `*Input`:

```go
type Command interface {
    Name() string
    Description() string
    Group() string
    Args() []Arg
    Options() []Option
    Execute(ctx context.Context, in io.Reader, out io.Writer, input *Input) error
}
```

### Аргументы

Позиционные параметры. Обязательные по умолчанию, порядок определяется порядком в слайсе:

```go
// Обязательный аргумент
cli.StringArg("direction", "Migration direction (up/down)")

// Необязательный аргумент со значением по умолчанию
cli.StringArgOptional("environment", "production", "Target environment")
```

> Обязательные аргументы должны идти перед необязательными — это проверяется при регистрации.

### Опции

Именованные флаги, всегда необязательны, всегда имеют значение по умолчанию. Поддерживаются три типа:

```go
cli.StringOption("format", "f", "json", "Output format")
cli.IntOption("timeout", "t", 30, "Timeout in seconds")
cli.BoolOption("verbose", "v", false, "Enable verbose output")
```

### `*Input`

Предоставляет доступ к распарсенным значениям:

```go
func (c *MyCommand) Execute(ctx context.Context, in io.Reader, out io.Writer, input *cli.Input) error {
    name := input.Arg("name")              // позиционный аргумент
    format := input.StringOption("format") // строковая опция
    count := input.IntOption("count")      // числовая опция
    verbose := input.BoolOption("verbose") // булева опция
    extra := input.RemainingArgs()         // неописанные позиционные аргументы
    // ...
}
```

### `Console`

Точка входа приложения. Настраивается через функциональные опции:

```go
app := cli.New(
    cli.WithName("myapp"),
    cli.WithVersion("1.0.0"),
    cli.WithMiddleware(loggingMiddleware, timingMiddleware),
)
```

### Middleware

Middleware оборачивает выполнение команды, позволяя добавлять сквозную логику. Применяется в порядке регистрации — первый зарегистрированный является внешним:

```go
type Handler func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error
type Middleware func(next Handler) Handler
```

### Типизированные ошибки

Все ошибки пакета типизированы и реализуют интерфейс `ExitCoder`:

| Ошибка | Код завершения | Когда возникает |
|--------|:-:|---|
| `CommandNotFoundError` | 2 | Команда не найдена в реестре |
| `MissingArgumentError` | 2 | Не передан обязательный аргумент |
| `CommandExistsError` | — | Повторная регистрация команды |
| `InvalidCommandError` | — | Невалидное определение команды |
| `PanicError` | 1 | Паника при выполнении команды |

```go
// Извлечение exit-кода из ошибки
code := cli.GetExitCode(err) // 0 для nil, 1 если ExitCoder не реализован
```

### Встроенные команды

- **`help`** — выводит список всех команд или справку по конкретной команде
- **`version`** — выводит версию приложения (регистрируется при указании `WithVersion`)

---

## 🧪 Примеры использования

### Минимальная команда

```go
package main

import (
    "context"
    "fmt"
    "io"
    "os"

    "github.com/shuldan/cli"
)

type GreetCommand struct{}

func (c *GreetCommand) Name() string        { return "greet" }
func (c *GreetCommand) Description() string { return "Print a greeting message" }
func (c *GreetCommand) Group() string       { return "misc" }

func (c *GreetCommand) Args() []cli.Arg {
    return []cli.Arg{
        cli.StringArg("name", "Name to greet"),
    }
}

func (c *GreetCommand) Options() []cli.Option {
    return []cli.Option{
        cli.StringOption("greeting", "g", "Hello", "Greeting word"),
        cli.BoolOption("loud", "", false, "UPPERCASE output"),
    }
}

func (c *GreetCommand) Execute(_ context.Context, _ io.Reader, out io.Writer, input *cli.Input) error {
    name := input.Arg("name")
    greeting := input.StringOption("greeting")

    message := fmt.Sprintf("%s, %s!", greeting, name)
    if input.BoolOption("loud") {
        message = strings.ToUpper(message)
    }

    _, err := fmt.Fprintln(out, message)
    return err
}

func main() {
    app := cli.New(
        cli.WithName("greeter"),
        cli.WithVersion("1.0.0"),
    )

    if err := app.Register(&GreetCommand{}); err != nil {
        log.Fatal(err)
    }

    ctx, cancel := cli.SignalContext(context.Background())
    defer cancel()

    if err := app.Run(ctx, os.Stdin, os.Stdout, os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(cli.GetExitCode(err))
    }
}
```

```sh
greeter greet Alice                     # Hello, Alice!
greeter greet Alice --greeting Hi       # Hi, Alice!
greeter greet Alice -g Hey --loud       # HEY, ALICE!
```

### Команда с зависимостями

```go
type MigrateCommand struct {
    db *sql.DB
}

func (c *MigrateCommand) Name() string        { return "migrate" }
func (c *MigrateCommand) Description() string { return "Run database migrations" }
func (c *MigrateCommand) Group() string       { return "database" }

func (c *MigrateCommand) Args() []cli.Arg {
    return []cli.Arg{
        cli.StringArg("direction", "Migration direction (up/down)"),
    }
}

func (c *MigrateCommand) Options() []cli.Option {
    return []cli.Option{
        cli.IntOption("steps", "s", 0, "Number of migration steps (0 = all)"),
        cli.BoolOption("dry-run", "", false, "Preview without applying"),
    }
}

func (c *MigrateCommand) Execute(ctx context.Context, _ io.Reader, out io.Writer, input *cli.Input) error {
    direction := input.Arg("direction")
    steps := input.IntOption("steps")
    dryRun := input.BoolOption("dry-run")

    fmt.Fprintf(out, "Migrating %s (steps=%d, dry-run=%v)\n", direction, steps, dryRun)
    // c.db используется для выполнения миграций
    return nil
}

// Регистрация — зависимости передаются через конструктор:
db, _ := sql.Open("postgres", dsn)
app.Register(&MigrateCommand{db: db})
```

```sh
myapp migrate up
myapp migrate down --steps=3
myapp migrate up --dry-run
```

### Middleware

```go
// Замер времени выполнения
func timingMiddleware(next cli.Handler) cli.Handler {
    return func(ctx context.Context, in io.Reader, out io.Writer, input *cli.Input) error {
        start := time.Now()
        err := next(ctx, in, out, input)
        log.Printf("executed in %s", time.Since(start))
        return err
    }
}

// Логирование ошибок
func errorLoggingMiddleware(next cli.Handler) cli.Handler {
    return func(ctx context.Context, in io.Reader, out io.Writer, input *cli.Input) error {
        err := next(ctx, in, out, input)
        if err != nil {
            log.Printf("command failed: %v", err)
        }
        return err
    }
}

app := cli.New(
    cli.WithName("myapp"),
    cli.WithMiddleware(errorLoggingMiddleware, timingMiddleware),
)
```

Порядок выполнения: `errorLogging → timing → команда → timing → errorLogging`.

### Справка

Справка генерируется автоматически из описаний команд, аргументов и опций:

```sh
# Общая справка
myapp help
myapp --help
myapp -h
myapp          # без аргументов — тоже справка

# Справка по конкретной команде
myapp help migrate
myapp migrate --help
```

Пример вывода `myapp help`:
```
myapp

Usage: myapp <command> [options] [arguments]

console:
  help       Display help for commands
  version    Display application version

database:
  migrate    Run database migrations
  seed       Seed database with test data

misc:
  greet      Print a greeting message
```

Пример вывода `myapp help migrate`:
```
migrate — Run database migrations

Usage: myapp migrate <direction> [--steps=...] [--dry-run]

Arguments:
  direction    Migration direction (up/down)

Options:
  --steps, -s    Number of migration steps (0 = all) (default: 0)
  --dry-run      Preview without applying (default: false)
```

### Версия

```sh
myapp version
myapp --version
myapp -v
```

```
myapp version 1.0.0
```

### Graceful shutdown

```go
func main() {
    // Контекст отменяется при SIGINT (Ctrl+C) или SIGTERM
    ctx, cancel := cli.SignalContext(context.Background())
    defer cancel()

    app := cli.New(cli.WithName("myapp"))
    // ...

    if err := app.Run(ctx, os.Stdin, os.Stdout, os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(cli.GetExitCode(err))
    }
}
```

### Корректная обработка exit-кодов

```go
if err := app.Run(ctx, os.Stdin, os.Stdout, os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(cli.GetExitCode(err))
    // CommandNotFoundError, MissingArgumentError → exit 2
    // PanicError, прочие ошибки                 → exit 1
    // nil                                       → exit 0
}
```

---

## 📐 Валидация при регистрации

Пакет проверяет корректность определения команды в момент вызова `Register`:

- Команда не `nil`, имя не пустое
- Имя команды уникально
- Обязательные аргументы идут перед необязательными
- Нет дублирующихся имён аргументов
- Нет дублирующихся имён и коротких имён опций

```go
err := app.Register(&BadCommand{})
// *cli.InvalidCommandError: invalid command "bad": required argument "y" cannot follow optional argument
// *cli.CommandExistsError: command already registered: migrate
```

---

## 📄 Лицензия

Проект распространяется под лицензией [MIT](LICENSE).

---

## 🤝 Вклад в проект

PR и issue приветствуются! Следуйте стандартам форматирования и покрывайте новый код тестами.

---

> **Автор**: MSeytumerov
> **Репозиторий**: `github.com/shuldan/cli`
> **Go**: `1.24.2`
