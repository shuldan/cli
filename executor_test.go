package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestExecutor_Execute_ContextCancelled(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	e := &executor{parser: newTestParser(reg), middleware: nil}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := e.execute(ctx, strings.NewReader(""), &bytes.Buffer{}, []string{cmdNameHelp})
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestExecutor_Execute_ParseError(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	e := &executor{parser: newTestParser(reg), middleware: nil}
	err := e.execute(
		context.Background(), strings.NewReader(""), &bytes.Buffer{}, []string{"nonexistent"},
	)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestExecutor_Execute_Success(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	called := false
	cmd := &mockCommand{
		name: "test",
		executeFn: func(_ context.Context, _ io.Reader, _ io.Writer, _ *Input) error {
			called = true
			return nil
		},
	}
	_ = reg.register(&helpCommand{registry: reg, appName: "app"})
	_ = reg.register(cmd)
	e := &executor{parser: newTestParser(reg), middleware: nil}
	err := e.execute(
		context.Background(), strings.NewReader(""), &bytes.Buffer{}, []string{"test"},
	)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Errorf("expected command to be called")
	}
}

func TestExecutor_Execute_WithMiddleware(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "test"}
	_ = reg.register(&helpCommand{registry: reg, appName: "app"})
	_ = reg.register(cmd)
	order := make([]string, 0, 2)
	mw1 := func(next Handler) Handler {
		return func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error {
			order = append(order, "mw1")
			return next(ctx, in, out, input)
		}
	}
	mw2 := func(next Handler) Handler {
		return func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error {
			order = append(order, "mw2")
			return next(ctx, in, out, input)
		}
	}
	e := &executor{parser: newTestParser(reg), middleware: []Middleware{mw1, mw2}}
	err := e.execute(
		context.Background(), strings.NewReader(""), &bytes.Buffer{}, []string{"test"},
	)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(order) != 2 || order[0] != "mw1" || order[1] != "mw2" {
		t.Errorf("expected [mw1 mw2], got %v", order)
	}
}

func TestExecutor_WithRecovery_Panic(t *testing.T) {
	t.Parallel()
	e := &executor{}
	err := e.withRecovery(func() error {
		panic("test panic")
	})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	pe, ok := err.(*PanicError)
	if !ok {
		t.Fatalf("expected *PanicError, got %T", err)
	}
	if pe.Value != "test panic" {
		t.Errorf("expected panic value %q, got %v", "test panic", pe.Value)
	}
}

func TestExecutor_WithRecovery_NoPanic(t *testing.T) {
	t.Parallel()
	e := &executor{}
	expected := fmt.Errorf("some error")
	err := e.withRecovery(func() error { return expected })
	if err != expected {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func TestExecutor_WithRecovery_NilReturn(t *testing.T) {
	t.Parallel()
	e := &executor{}
	err := e.withRecovery(func() error { return nil })
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestExecutor_Execute_CommandError(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmdErr := fmt.Errorf("command failed")
	cmd := &mockCommand{
		name: "failing",
		executeFn: func(_ context.Context, _ io.Reader, _ io.Writer, _ *Input) error {
			return cmdErr
		},
	}
	_ = reg.register(&helpCommand{registry: reg, appName: "app"})
	_ = reg.register(cmd)
	e := &executor{parser: newTestParser(reg), middleware: nil}
	err := e.execute(
		context.Background(), strings.NewReader(""), &bytes.Buffer{}, []string{"failing"},
	)
	if err != cmdErr {
		t.Errorf("expected %v, got %v", cmdErr, err)
	}
}
