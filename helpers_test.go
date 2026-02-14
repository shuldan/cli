package cli

import (
	"context"
	"fmt"
	"io"
)

type mockCommand struct {
	name        string
	description string
	group       string
	args        []Arg
	options     []Option
	executeFn   func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error
}

func (m *mockCommand) Name() string        { return m.name }
func (m *mockCommand) Description() string { return m.description }
func (m *mockCommand) Group() string       { return m.group }
func (m *mockCommand) Args() []Arg         { return m.args }
func (m *mockCommand) Options() []Option   { return m.options }

func (m *mockCommand) Execute(ctx context.Context, in io.Reader, out io.Writer, input *Input) error {
	if m.executeFn != nil {
		return m.executeFn(ctx, in, out, input)
	}
	return nil
}

type failAfterNWriter struct {
	n       int
	written int
}

func (w *failAfterNWriter) Write(p []byte) (int, error) {
	w.written++
	if w.written > w.n {
		return 0, fmt.Errorf("write error")
	}
	return len(p), nil
}

func newTestRegistry() *registry {
	return &registry{
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
}

func newTestParser(reg *registry) *parser {
	return &parser{registry: reg}
}
