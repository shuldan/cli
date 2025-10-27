package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

type mockCommand struct {
	name        string
	description string
	group       string
	shouldPanic bool
	executeErr  error
}

func (m *mockCommand) Name() string            { return m.name }
func (m *mockCommand) Description() string     { return m.description }
func (m *mockCommand) Group() string           { return m.group }
func (m *mockCommand) Configure(*flag.FlagSet) {}
func (m *mockCommand) Execute(_ context.Context, _ io.Reader, _ io.Writer, _ []string) error {
	if m.shouldPanic {
		panic("test panic")
	}
	return m.executeErr
}

func TestConsole_Register_Success(t *testing.T) {
	console := New()
	cmd := &mockCommand{name: "test", description: "Test command", group: "test"}

	err := console.Register(cmd)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestConsole_Register_ErrorCases(t *testing.T) {
	console := New()

	t.Run("NilCommand", func(t *testing.T) {
		err := console.Register(nil)
		if err == nil {
			t.Error("Expected error for nil command, got nil")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		cmd := &mockCommand{name: ""}
		err := console.Register(cmd)
		if err == nil {
			t.Error("Expected error for empty name, got nil")
		}
	})

	t.Run("DuplicateCommand", func(t *testing.T) {
		cmd1 := &mockCommand{name: "duplicate", description: "First"}
		cmd2 := &mockCommand{name: "duplicate", description: "Second"}

		_ = console.Register(cmd1)
		err := console.Register(cmd2)
		if err == nil {
			t.Error("Expected error for duplicate command, got nil")
		}
	})
}

func TestConsole_Run_Success(t *testing.T) {
	console := New()
	cmd := &mockCommand{name: "test", description: "Test command", group: "test"}
	_ = console.Register(cmd)

	ctx := context.Background()
	input := strings.NewReader("")
	output := &strings.Builder{}
	args := []string{"test"}

	err := console.Run(ctx, input, output, args)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestConsole_Run_ErrorCases(t *testing.T) {
	console := New()

	t.Run("NoArguments", func(t *testing.T) {
		ctx := context.Background()
		input := strings.NewReader("")
		output := &strings.Builder{}
		var args []string

		err := console.Run(ctx, input, output, args)
		if err == nil {
			t.Error("Expected error for no arguments, got nil")
		}
	})

	t.Run("UnknownCommand", func(t *testing.T) {
		ctx := context.Background()
		input := strings.NewReader("")
		output := &strings.Builder{}
		args := []string{"unknown"}

		err := console.Run(ctx, input, output, args)
		if err == nil {
			t.Error("Expected error for unknown command, got nil")
		}
	})

	t.Run("CommandPanic", func(t *testing.T) {
		cmd := &mockCommand{name: "panic", description: "Panic command", group: "test", shouldPanic: true}
		_ = console.Register(cmd)

		ctx := context.Background()
		input := strings.NewReader("")
		output := &strings.Builder{}
		args := []string{"panic"}

		err := console.Run(ctx, input, output, args)
		if err == nil {
			t.Error("Expected error for panic, got nil")
		}
		if !strings.Contains(err.Error(), "panic during command execution") {
			t.Errorf("Expected panic error, got %v", err)
		}
	})
}

func TestRegistry_Register_Success(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	cmd := &mockCommand{name: "test", description: "Test command", group: "test"}
	err := reg.register(cmd)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(reg.commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(reg.commands))
	}

	if len(reg.groups["test"]) != 1 {
		t.Errorf("Expected 1 command in group, got %d", len(reg.groups["test"]))
	}
}

func TestRegistry_Register_ErrorCases(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	t.Run("NilCommand", func(t *testing.T) {
		err := reg.register(nil)
		if err == nil {
			t.Error("Expected error for nil command, got nil")
		}
	})

	t.Run("EmptyName", func(t *testing.T) {
		cmd := &mockCommand{name: ""}
		err := reg.register(cmd)
		if err == nil {
			t.Error("Expected error for empty name, got nil")
		}
	})

	t.Run("DuplicateCommand", func(t *testing.T) {
		cmd1 := &mockCommand{name: "duplicate", description: "First"}
		cmd2 := &mockCommand{name: "duplicate", description: "Second"}

		_ = reg.register(cmd1)
		err := reg.register(cmd2)
		if err == nil {
			t.Error("Expected error for duplicate command, got nil")
		}
	})
}

func TestParser_Parse_Success(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	parser := &parser{registry: reg}

	cmd := &mockCommand{name: "test", description: "Test command", group: "test"}
	_ = reg.register(cmd)

	output := &strings.Builder{}
	args := []string{"test"}

	parsed, err := parser.parse(output, args)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if parsed == nil {
		t.Fatal("Expected parsed command, got nil")
	}
	if parsed.Name != "test" {
		t.Errorf("Expected command name 'test', got %s", parsed.Name)
	}
}

func TestParser_Parse_ErrorCases(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	parser := &parser{registry: reg}

	t.Run("NoArguments", func(t *testing.T) {
		output := &strings.Builder{}
		var args []string

		_, err := parser.parse(output, args)
		if err == nil {
			t.Error("Expected error for no arguments, got nil")
		}
	})

	t.Run("UnknownCommand", func(t *testing.T) {
		output := &strings.Builder{}
		args := []string{"unknown"}

		_, err := parser.parse(output, args)
		if err == nil {
			t.Error("Expected error for unknown command, got nil")
		}
	})
}

func TestExecutor_Execute_Success(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	exec := &executor{
		parser: &parser{registry: reg},
	}

	cmd := &mockCommand{name: "test", description: "Test command", group: "test"}
	_ = reg.register(cmd)

	ctx := context.Background()
	input := strings.NewReader("")
	output := &strings.Builder{}
	args := []string{"test"}

	err := exec.execute(ctx, input, output, args)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExecutor_Execute_ErrorCases(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	t.Run("ParseError", func(t *testing.T) {
		exec := &executor{
			parser: &parser{registry: reg},
		}

		ctx := context.Background()
		input := strings.NewReader("")
		output := &strings.Builder{}
		var args []string

		err := exec.execute(ctx, input, output, args)
		if err == nil {
			t.Error("Expected error for parse error, got nil")
		}
	})

	t.Run("CommandPanic", func(t *testing.T) {
		exec := &executor{
			parser: &parser{registry: reg},
		}

		cmd := &mockCommand{name: "panic", description: "Panic command", group: "test", shouldPanic: true}
		_ = reg.register(cmd)

		ctx := context.Background()
		input := strings.NewReader("")
		output := &strings.Builder{}
		args := []string{"panic"}

		err := exec.execute(ctx, input, output, args)
		if err == nil {
			t.Error("Expected error for panic, got nil")
		}
		if !strings.Contains(err.Error(), "panic during command execution") {
			t.Errorf("Expected panic error, got %v", err)
		}
	})
}

func TestHelp_Execute_General(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	helpCmd := &help{registry: reg}

	cmd1 := &mockCommand{name: "cmd1", description: "Command 1", group: "group1"}
	cmd2 := &mockCommand{name: "cmd2", description: "Command 2", group: "group2"}
	_ = reg.register(cmd1)
	_ = reg.register(cmd2)

	ctx := context.Background()
	input := strings.NewReader("")
	output := &strings.Builder{}
	var args []string

	err := helpCmd.Execute(ctx, input, output, args)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "cmd1") || !strings.Contains(result, "cmd2") {
		t.Error("Expected commands in output")
	}
	if !strings.Contains(result, "group1") || !strings.Contains(result, "group2") {
		t.Error("Expected groups in output")
	}
}

func TestHelp_Execute_SpecificCommand(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	helpCmd := &help{registry: reg, command: "test"}

	cmd := &mockCommand{name: "test", description: "Test command", group: "test"}
	_ = reg.register(cmd)

	ctx := context.Background()
	input := strings.NewReader("")
	output := &strings.Builder{}
	var args []string

	err := helpCmd.Execute(ctx, input, output, args)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	result := output.String()
	if !strings.Contains(result, "test - Test command") {
		t.Errorf("Expected command help in output, got %s", result)
	}
}

func TestHelp_Execute_CommandNotFound(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	helpCmd := &help{registry: reg, command: "unknown"}

	ctx := context.Background()
	input := strings.NewReader("")
	output := &strings.Builder{}
	var args []string

	err := helpCmd.Execute(ctx, input, output, args)
	if err == nil {
		t.Error("Expected error for unknown command, got nil")
	}
}

func TestCommand_Interface(t *testing.T) {
	cmd := &mockCommand{
		name:        "test",
		description: "Test command",
		group:       "test",
	}

	if cmd.Name() != "test" {
		t.Errorf("Expected name 'test', got %s", cmd.Name())
	}

	if cmd.Description() != "Test command" {
		t.Errorf("Expected description 'Test command', got %s", cmd.Description())
	}

	if cmd.Group() != "test" {
		t.Errorf("Expected group 'test', got %s", cmd.Group())
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Configure panicked: %v", r)
			}
		}()
		cmd.Configure(&flag.FlagSet{})
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Execute panicked: %v", r)
			}
		}()
		err := cmd.Execute(context.Background(), nil, nil, nil)
		if err != nil {
			t.Errorf("Expected no error from Execute, got %v", err)
		}
	}()
}

func TestHelp_Description(t *testing.T) {
	h := &help{}
	expected := "Display help for commands"
	if h.Description() != expected {
		t.Errorf("Expected description %q, got %q", expected, h.Description())
	}
}

func TestHelp_Configure(t *testing.T) {
	h := &help{}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	h.Configure(fs)

	fs.VisitAll(func(f *flag.Flag) {
		if f.Name != "command" {
			t.Errorf("Expected flag 'command', got %s", f.Name)
		}
	})
}

func TestHelp_Execute_GeneralHelpError(t *testing.T) {

	failWriter := &failingWriter{}

	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	helpCmd := &help{registry: reg}

	ctx := context.Background()
	input := strings.NewReader("")
	var args []string

	err := helpCmd.Execute(ctx, input, failWriter, args)
	if err == nil {
		t.Error("Expected error from failing writer, got nil")
	}
}

func TestHelp_ShowCommandHelp_ErrorCases(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	helpCmd := &help{registry: reg}

	t.Run("FailingNameDescriptionWrite", func(t *testing.T) {
		cmd := &mockCommand{name: "test", description: "Test command", group: "test"}
		_ = reg.register(cmd)

		failWriter := &failingWriter{}
		err := helpCmd.showCommandHelp("test", failWriter)
		if err == nil {
			t.Error("Expected error from failing writer, got nil")
		}
	})

	t.Run("FailingOptionsWrite", func(t *testing.T) {
		cmd := &mockCommandWithFlags{name: "test", description: "Test command"}
		_ = reg.register(cmd)

		failWriter := &partialFailingWriter{failOnCall: 2}
		err := helpCmd.showCommandHelp("test", failWriter)
		if err == nil {
			t.Error("Expected error from failing writer, got nil")
		}
	})
}

func TestParser_Parse_FlagParsingError(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	parser := &parser{registry: reg}

	cmd := &mockCommandWithFlags{name: "test", description: "Test command"}
	_ = reg.register(cmd)

	output := &strings.Builder{}

	args := []string{"test", "-invalid"}

	_, err := parser.parse(output, args)
	if err == nil {
		t.Error("Expected error for invalid flag, got nil")
	}
}

func TestRegistry_GetGroups_Sorting(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	cmdC := &mockCommand{name: "c_command", group: "test"}
	cmdA := &mockCommand{name: "a_command", group: "test"}
	cmdB := &mockCommand{name: "b_command", group: "test"}

	_ = reg.register(cmdC)
	_ = reg.register(cmdA)
	_ = reg.register(cmdB)

	groups := reg.getGroups()
	commands := groups["test"]

	if len(commands) != 3 {
		t.Fatalf("Expected 3 commands, got %d", len(commands))
	}

	expectedOrder := []string{"a_command", "b_command", "c_command"}
	for i, expectedName := range expectedOrder {
		if commands[i].Name() != expectedName {
			t.Errorf("Expected command %s at position %d, got %s", expectedName, i, commands[i].Name())
		}
	}
}

func TestRegistry_GetGroups_EmptyGroup(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}

	reg.groups["empty"] = []string{"nonexistent"}

	groups := reg.getGroups()
	commands := groups["empty"]

	if len(commands) != 0 {
		t.Errorf("Expected 0 commands in empty group, got %d", len(commands))
	}
}

type failingWriter struct{}

func (fw *failingWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("write failed")
}

type partialFailingWriter struct {
	callCount  int
	failOnCall int
}

func (pw *partialFailingWriter) Write(p []byte) (int, error) {
	pw.callCount++
	if pw.callCount == pw.failOnCall {
		return 0, fmt.Errorf("write failed on call %d", pw.callCount)
	}
	return len(p), nil
}

type mockCommandWithFlags struct {
	name        string
	description string
	flagValue   string
}

func (m *mockCommandWithFlags) Name() string        { return m.name }
func (m *mockCommandWithFlags) Description() string { return m.description }
func (m *mockCommandWithFlags) Group() string       { return "test" }
func (m *mockCommandWithFlags) Configure(fs *flag.FlagSet) {
	fs.StringVar(&m.flagValue, "testflag", "", "Test flag")
}
func (m *mockCommandWithFlags) Execute(_ context.Context, _ io.Reader, _ io.Writer, _ []string) error {
	return nil
}

func TestExecutor_Execute_CancelledContext(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	exec := &executor{
		parser: &parser{registry: reg},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	input := strings.NewReader("")
	output := &strings.Builder{}
	args := []string{"test"}

	err := exec.execute(ctx, input, output, args)
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

func TestExecutor_Execute_DeadlineExceededContext(t *testing.T) {
	reg := &registry{
		mu:       sync.RWMutex{},
		commands: make(map[string]Command),
		groups:   make(map[string][]string),
	}
	exec := &executor{
		parser: &parser{registry: reg},
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	input := strings.NewReader("")
	output := &strings.Builder{}
	args := []string{"test"}

	err := exec.execute(ctx, input, output, args)
	if err == nil {
		t.Error("Expected error for expired deadline, got nil")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}
}
