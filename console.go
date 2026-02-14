package cli

import (
	"context"
	"fmt"
	"io"
	"sync"
)

type Console struct {
	name       string
	version    string
	registry   *registry
	executor   *executor
	middleware []Middleware
}

type ConsoleOption func(*Console)

func WithName(name string) ConsoleOption {
	return func(c *Console) {
		c.name = name
	}
}

func WithVersion(version string) ConsoleOption {
	return func(c *Console) {
		c.version = version
	}
}

func WithMiddleware(m ...Middleware) ConsoleOption {
	return func(c *Console) {
		c.middleware = append(c.middleware, m...)
	}
}

func New(opts ...ConsoleOption) *Console {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	c := &Console{
		name:     "app",
		registry: reg,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.executor = &executor{
		parser: &parser{
			registry: reg,
		},
		middleware: c.middleware,
	}

	c.mustRegister(&helpCommand{
		registry: reg,
		appName:  c.name,
	})

	if c.version != "" {
		c.mustRegister(&versionCommand{
			appName:    c.name,
			appVersion: c.version,
		})
	}

	return c
}

func (c *Console) Register(cmd Command) error {
	return c.registry.register(cmd)
}

func (c *Console) Run(ctx context.Context, in io.Reader, out io.Writer, args []string) error {
	return c.executor.execute(ctx, in, out, args)
}

func (c *Console) mustRegister(cmd Command) {
	if err := c.registry.register(cmd); err != nil {
		panic(fmt.Sprintf("cli: failed to register built-in command %q: %v", cmd.Name(), err))
	}
}
