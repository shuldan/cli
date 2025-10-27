package console

import (
	"context"
	"io"
	"sync"
)

type Console struct {
	registry *registry
	executor *executor
}

func New() *Console {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	exec := &executor{
		parser: &parser{
			registry: reg,
		},
	}

	cli := &Console{
		registry: reg,
		executor: exec,
	}

	_ = cli.Register(&help{
		registry: reg,
	})

	return cli
}

func (c *Console) Register(cmd Command) error {
	return c.registry.register(cmd)
}

func (c *Console) Run(ctx context.Context, input io.Reader, output io.Writer, args []string) error {
	return c.executor.execute(ctx, input, output, args)
}
