package cli

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	t.Parallel()
	c := New()
	if c.name != "app" {
		t.Errorf("expected name %q, got %q", "app", c.name)
	}
	if c.version != "" {
		t.Errorf("expected empty version, got %q", c.version)
	}
	if _, ok := c.registry.get("help"); !ok {
		t.Errorf("expected help command to be registered")
	}
	if _, ok := c.registry.get("version"); ok {
		t.Errorf("expected version command NOT to be registered")
	}
}

func TestNew_WithOptions(t *testing.T) {
	t.Parallel()
	mw := func(next Handler) Handler { return next }
	c := New(WithName("myapp"), WithVersion("1.0"), WithMiddleware(mw))
	if c.name != "myapp" {
		t.Errorf("expected name %q, got %q", "myapp", c.name)
	}
	if c.version != "1.0" {
		t.Errorf("expected version %q, got %q", "1.0", c.version)
	}
	if _, ok := c.registry.get("version"); !ok {
		t.Errorf("expected version command to be registered")
	}
	if len(c.middleware) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(c.middleware))
	}
}

func TestConsole_Register_Success(t *testing.T) {
	t.Parallel()
	c := New()
	cmd := &mockCommand{name: "test", description: "test cmd"}
	if err := c.Register(cmd); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestConsole_Register_Duplicate(t *testing.T) {
	t.Parallel()
	c := New()
	cmd := &mockCommand{name: "test", description: "test cmd"}
	_ = c.Register(cmd)
	err := c.Register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestConsole_Run_Success(t *testing.T) {
	t.Parallel()
	c := New()
	called := false
	cmd := &mockCommand{
		name: "greet",
		executeFn: func(_ context.Context, _ io.Reader, _ io.Writer, _ *Input) error {
			called = true
			return nil
		},
	}
	_ = c.Register(cmd)
	var buf bytes.Buffer
	err := c.Run(context.Background(), strings.NewReader(""), &buf, []string{"greet"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !called {
		t.Errorf("expected command to be called")
	}
}

func TestConsole_mustRegister_Panic(t *testing.T) {
	t.Parallel()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	c := New()
	c.mustRegister(nil)
}
