# `console` — Минималистичный CLI-фреймворк для Go-приложений

[![Go CI](https://github.com/shuldan/app/workflows/Go%20CI/badge.svg)](https://github.com/shuldan/app/actions)  
[![codecov](https://codecov.io/gh/shuldan/app/branch/main/graph/badge.svg)](https://codecov.io/gh/shuldan/app)  
[![Go Report Card](https://goreportcard.com/badge/github.com/shuldan/app)](https://goreportcard.com/report/github.com/shuldan/app)  
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Пакет `console` предоставляет простой, но расширяемый каркас для создания консольных утилит на Go с поддержкой регистрации команд, автоматической генерацией справки и безопасным выполнением.

---

## 🚀 Основные возможности

- **Регистрация команд**: каждая команда реализует единый интерфейс `Command`.
- **Автоматическая справка**: встроенная команда `help` выводит список всех команд и детали по конкретной.
- **Группировка команд**: команды можно объединять в логические группы (например, `database`, `cache`, `console`).
- **Безопасное выполнение**: перехват паник внутри команд с корректным отчётом об ошибке.
- **Поддержка флагов**: каждая команда может определять собственные флаги через стандартный `flag.FlagSet`.
- **Тестируемость**: полная поддержка внедрения зависимостей (`io.Reader`, `io.Writer`, `context.Context`).

---

## 📦 Установка зависимостей и инструментов

Для работы с проектом требуется Go 1.24+.

Установите необходимые инструменты:

```sh
make install-tools
```

Это установит:
- `golangci-lint` (v2.4.0)
- `goimports`
- `gosec`

---

## 🛠️ Работа с проектом

### Запуск локальной проверки

```sh
make all
```

Выполняет:
- проверку форматирования кода,
- статический анализ (`golangci-lint`),
- security-сканирование (`gosec`),
- запуск тестов.

### Проверка в CI

```sh
make ci
```

Запускает:
- форматирование,
- `go vet`,
- линтер,
- тесты с отчётом о покрытии.

### Форматирование кода

```sh
make fmt
```

Автоматически форматирует `.go` файлы и сортирует импорты.

### Запуск тестов

```sh
make test
make test-coverage
```

---

## 🧱 Архитектура

### `Console`

Основной объект CLI-приложения:

```go
cli := console.New()
err := cli.Register(&MyCommand{})
err = cli.Run(ctx, os.Stdin, os.Stdout, os.Args[1:])
```

### `Command`

Любая команда должна реализовывать интерфейс:

```go
type Command interface {
	Name() string
	Description() string
	Group() string
	Configure(flags *flag.FlagSet)
	Execute(ctx context.Context, input io.Reader, output io.Writer, args []string) error
}
```

Команда `help` регистрируется автоматически при создании `Console`.

---

## 🧪 Пример использования

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shuldan/app/cli/pkg/console"
)

type GreetCommand struct {
	name string
}

func (g *GreetCommand) Name() string        { return "greet" }
func (g *GreetCommand) Description() string { return "Print greeting message" }
func (g *GreetCommand) Group() string       { return "misc" }

func (g *GreetCommand) Configure(flags *flag.FlagSet) {
	flags.StringVar(&g.name, "name", "World", "Name to greet")
}

func (g *GreetCommand) Execute(_ context.Context, _ io.Reader, w io.Writer, _ []string) error {
	_, err := fmt.Fprintf(w, "Hello, %s!\n", g.name)
	return err
}

func main() {
	cli := console.New()
	if err := cli.Register(&GreetCommand{}); err != nil {
		panic(err)
	}

	ctx := context.Background()
	if err := cli.Run(ctx, os.Stdin, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

Примеры вызова:
```sh
./mycli greet --name Alice
./mycli help
./mycli help --command greet
```

---

## 📄 Лицензия

Проект распространяется под лицензией [MIT](LICENSE).

---

## 🤝 Вклад в проект

PR и issue приветствуются! Следуйте стандартам форматирования и покрывайте новый код тестами.

---

> **Автор**: MSeytumerov  
> **Репозиторий**: `github.com/shuldan/app`  
> **Пакет**: `github.com/shuldan/app/cli/pkg/console`  
> **Go**: `1.24.2`