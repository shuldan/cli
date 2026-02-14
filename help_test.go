package cli

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
)

func newHelpCmd() (*helpCommand, *registry) {
	reg := newTestRegistry()
	h := &helpCommand{registry: reg, appName: "testapp"}
	_ = reg.register(h)
	return h, reg
}

func TestHelpCommand_Metadata(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	if h.Name() != "help" {
		t.Errorf("expected help, got %s", h.Name())
	}
	if h.Description() == "" {
		t.Errorf("expected non-empty description")
	}
	if h.Group() != "console" {
		t.Errorf("expected console, got %s", h.Group())
	}
	if len(h.Args()) != 1 {
		t.Errorf("expected 1 arg, got %d", len(h.Args()))
	}
	if h.Options() != nil {
		t.Errorf("expected nil options")
	}
}

func TestHelpCommand_Execute_GeneralHelp(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	inp := newInput(map[string]string{"command": ""}, nil, nil)
	err := h.Execute(context.Background(), nil, &buf, inp)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), "testapp") {
		t.Errorf("expected output to contain 'testapp'")
	}
}

func TestHelpCommand_Execute_CommandHelp(t *testing.T) {
	t.Parallel()
	h, reg := newHelpCmd()
	_ = reg.register(&mockCommand{name: "greet", description: "Greet"})
	var buf bytes.Buffer
	inp := newInput(map[string]string{"command": "greet"}, nil, nil)
	err := h.Execute(context.Background(), nil, &buf, inp)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), "greet") {
		t.Errorf("expected output to contain 'greet'")
	}
}

func TestHelpCommand_Execute_CommandNotFound(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	inp := newInput(map[string]string{"command": "missing"}, nil, nil)
	err := h.Execute(context.Background(), nil, &buf, inp)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_BuildUsageLine_AllTypes(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{
			{Name: "req"},
			{Name: "opt", Optional: true},
		},
		options: []Option{
			{Name: "flag", Type: OptionTypeBool},
			{Name: "out", Type: OptionTypeString},
		},
	}
	usage := h.buildUsageLine(cmd)
	if !strings.Contains(usage, "<req>") {
		t.Errorf("expected <req> in usage, got %s", usage)
	}
	if !strings.Contains(usage, "[opt]") {
		t.Errorf("expected [opt] in usage, got %s", usage)
	}
	if !strings.Contains(usage, "[--flag]") {
		t.Errorf("expected [--flag] in usage, got %s", usage)
	}
	if !strings.Contains(usage, "[--out=...]") {
		t.Errorf("expected [--out=...] in usage, got %s", usage)
	}
}

func TestHelpCommand_PrintArgs_Empty(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	err := h.printArgs(&buf, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output")
	}
}

func TestHelpCommand_PrintArgs_WithDefault(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	args := []Arg{{Name: "a", Optional: true, Default: "x", Description: "desc"}}
	err := h.printArgs(&buf, args)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !strings.Contains(buf.String(), "default: x") {
		t.Errorf("expected default output, got %q", buf.String())
	}
}

func TestHelpCommand_PrintOptions_Empty(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	err := h.printOptions(&buf, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output")
	}
}

func TestHelpCommand_PrintOptions_WithShortAndDefault(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	opts := []Option{
		{Name: "verbose", Short: "v", Description: "desc", Default: true},
	}
	err := h.printOptions(&buf, opts)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "--verbose") {
		t.Errorf("expected --verbose")
	}
	if !strings.Contains(out, "-v") {
		t.Errorf("expected -v")
	}
}

func TestHelpCommand_PrintOptions_NoShortNoDefault(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	var buf bytes.Buffer
	opts := []Option{{Name: "x", Description: "desc", Default: nil}}
	err := h.printOptions(&buf, opts)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	out := buf.String()
	if strings.Contains(out, ", -") {
		t.Errorf("did not expect short option marker")
	}
}

func TestHelpCommand_ShowGeneralHelp_WriteError(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	w := &failAfterNWriter{n: 0}
	err := h.showGeneralHelp(w)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_ShowCommandHelp_WriteError(t *testing.T) {
	t.Parallel()
	h, reg := newHelpCmd()
	_ = reg.register(&mockCommand{name: "cmd", description: "d"})
	w := &failAfterNWriter{n: 0}
	err := h.showCommandHelp("cmd", w)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_PrintArgs_WriteError(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	w := &failAfterNWriter{n: 0}
	args := []Arg{{Name: "a", Description: "d"}}
	err := h.printArgs(w, args)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_PrintArgs_WriteErrorOnEntry(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	w := &failAfterNWriter{n: 1}
	args := []Arg{{Name: "a", Description: "d"}}
	err := h.printArgs(w, args)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_PrintOptions_WriteError(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	w := &failAfterNWriter{n: 0}
	opts := []Option{{Name: "x", Description: "d"}}
	err := h.printOptions(w, opts)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_PrintOptions_WriteErrorOnEntry(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	w := &failAfterNWriter{n: 1}
	opts := []Option{{Name: "x", Description: "d"}}
	err := h.printOptions(w, opts)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestHelpCommand_ShowGeneralHelp_GroupWriteErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		failAt int
	}{
		{"fail at group header", 1},
		{"fail at command entry", 2},
		{"fail at blank line", 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			reg := newTestRegistry()
			h := &helpCommand{registry: reg, appName: "app"}
			_ = reg.register(h)
			w := &failAfterNWriter{n: tc.failAt}
			err := h.showGeneralHelp(w)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestHelpCommand_ShowCommandHelp_WithArgsAndOptions(t *testing.T) {
	t.Parallel()
	h, reg := newHelpCmd()
	cmd := &mockCommand{
		name:        "full",
		description: "Full cmd",
		args:        []Arg{{Name: "a", Description: "arg"}},
		options:     []Option{{Name: "o", Description: "opt", Default: "d"}},
	}
	_ = reg.register(cmd)
	var buf bytes.Buffer
	err := h.showCommandHelp("full", &buf)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Arguments") {
		t.Errorf("expected Arguments section")
	}
	if !strings.Contains(out, "Options") {
		t.Errorf("expected Options section")
	}
}

func TestHelpCommand_Execute_WriterInterface(t *testing.T) {
	t.Parallel()
	h, _ := newHelpCmd()
	inp := newInput(map[string]string{}, nil, nil)
	err := h.Execute(context.Background(), nil, io.Discard, inp)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
