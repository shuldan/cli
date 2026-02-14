package cli

import (
	"bytes"
	"strings"
	"testing"
)

func setupParserWithHelp() (*parser, *registry) {
	reg := newTestRegistry()
	_ = reg.register(&helpCommand{registry: reg, appName: "app"})
	return newTestParser(reg), reg
}

func TestParser_Parse_EmptyArgs(t *testing.T) {
	t.Parallel()
	p, _ := setupParserWithHelp()
	result, err := p.parse(&bytes.Buffer{}, []string{})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Command.Name() != "help" {
		t.Errorf("expected help, got %s", result.Command.Name())
	}
}

func TestParser_Parse_HelpFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
	}{
		{"--help", []string{"--help"}},
		{"-h", []string{"-h"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, _ := setupParserWithHelp()
			result, err := p.parse(&bytes.Buffer{}, tc.args)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result.Command.Name() != "help" {
				t.Errorf("expected help, got %s", result.Command.Name())
			}
		})
	}
}

func TestParser_Parse_VersionFlags(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args []string
	}{
		{"--version", []string{"--version"}},
		{"-v", []string{"-v"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p, reg := setupParserWithHelp()
			_ = reg.register(&versionCommand{appName: "app", appVersion: "1.0"})
			result, err := p.parse(&bytes.Buffer{}, tc.args)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result.Command.Name() != "version" {
				t.Errorf("expected version, got %s", result.Command.Name())
			}
		})
	}
}

func TestParser_Parse_VersionNotRegistered(t *testing.T) {
	t.Parallel()
	p, _ := setupParserWithHelp()
	_, err := p.parse(&bytes.Buffer{}, []string{"--version"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_Parse_CommandNotFound(t *testing.T) {
	t.Parallel()
	p, _ := setupParserWithHelp()
	_, err := p.parse(&bytes.Buffer{}, []string{"nonexistent"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_Parse_CommandWithArgs(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name: "greet",
		args: []Arg{StringArg("name", "Name")},
	}
	_ = reg.register(cmd)
	result, err := p.parse(&bytes.Buffer{}, []string{"greet", "Alice"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Input.Arg("name") != "Alice" {
		t.Errorf("expected Alice, got %q", result.Input.Arg("name"))
	}
}

func TestParser_Parse_WithOptions(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name: "cmd",
		options: []Option{
			StringOption("output", "o", "def", "Output"),
			IntOption("count", "c", 0, "Count"),
			BoolOption("verbose", "V", false, "Verbose"),
		},
	}
	_ = reg.register(cmd)
	result, err := p.parse(&bytes.Buffer{}, []string{"cmd", "--output=file", "-c", "5", "-V"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Input.StringOption("output") != "file" {
		t.Errorf("expected file, got %q", result.Input.StringOption("output"))
	}
	if result.Input.IntOption("count") != 5 {
		t.Errorf("expected 5, got %d", result.Input.IntOption("count"))
	}
	if !result.Input.BoolOption("verbose") {
		t.Errorf("expected verbose=true")
	}
}

func TestParser_Parse_MissingRequiredArg(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{StringArg("required", "Required arg")},
	}
	_ = reg.register(cmd)
	_, err := p.parse(&bytes.Buffer{}, []string{"cmd"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_Parse_RemainingArgs(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{StringArg("first", "First")},
	}
	_ = reg.register(cmd)
	result, err := p.parse(&bytes.Buffer{}, []string{"cmd", "a", "b", "c"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	rem := result.Input.RemainingArgs()
	if len(rem) != 2 || rem[0] != "b" || rem[1] != "c" {
		t.Errorf("expected [b c], got %v", rem)
	}
}

func TestParser_Parse_OptionalArgDefault(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{StringArgOptional("opt", "defval", "Optional")},
	}
	_ = reg.register(cmd)
	result, err := p.parse(&bytes.Buffer{}, []string{"cmd"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Input.Arg("opt") != "defval" {
		t.Errorf("expected defval, got %q", result.Input.Arg("opt"))
	}
}

func TestParser_BindOptions_InvalidDefaultString(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name:    "cmd",
		options: []Option{{Name: "x", Type: OptionTypeString, Default: 123}},
	}
	_ = reg.register(cmd)
	_, err := p.parse(&bytes.Buffer{}, []string{"cmd"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not a string") {
		t.Errorf("expected 'not a string', got %q", err.Error())
	}
}

func TestParser_BindOptions_InvalidDefaultInt(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name:    "cmd",
		options: []Option{{Name: "x", Type: OptionTypeInt, Default: "bad"}},
	}
	_ = reg.register(cmd)
	_, err := p.parse(&bytes.Buffer{}, []string{"cmd"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_BindOptions_InvalidDefaultBool(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name:    "cmd",
		options: []Option{{Name: "x", Type: OptionTypeBool, Default: "bad"}},
	}
	_ = reg.register(cmd)
	_, err := p.parse(&bytes.Buffer{}, []string{"cmd"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_BindOptions_UnsupportedType(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{
		name:    "cmd",
		options: []Option{{Name: "x", Type: OptionType(99), Default: nil}},
	}
	_ = reg.register(cmd)
	_, err := p.parse(&bytes.Buffer{}, []string{"cmd"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_CollectOptions_WrongPointerTypes(t *testing.T) {
	t.Parallel()
	p, _ := setupParserWithHelp()
	options := []Option{
		{Name: "s", Type: OptionTypeString},
		{Name: "i", Type: OptionTypeInt},
		{Name: "b", Type: OptionTypeBool},
	}
	pointers := map[string]any{
		"s": 123,
		"i": "wrong",
		"b": 42,
	}
	result := p.collectOptions(options, pointers)
	if _, ok := result["s"]; ok {
		t.Errorf("expected no string result for wrong pointer type")
	}
	if _, ok := result["i"]; ok {
		t.Errorf("expected no int result for wrong pointer type")
	}
	if _, ok := result["b"]; ok {
		t.Errorf("expected no bool result for wrong pointer type")
	}
}

func TestParser_BuildHelpParsed_NoHelpRegistered(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	p := newTestParser(reg)
	_, err := p.buildHelpParsed("")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_BuildHelpParsed_WithCommand(t *testing.T) {
	t.Parallel()
	p, _ := setupParserWithHelp()
	result, err := p.buildHelpParsed("test")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Input.Arg("command") != "test" {
		t.Errorf("expected test, got %q", result.Input.Arg("command"))
	}
}

func TestParser_BuildHelpParsed_Empty(t *testing.T) {
	t.Parallel()
	p, _ := setupParserWithHelp()
	result, err := p.buildHelpParsed("")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Input.Arg("command") != "" {
		t.Errorf("expected empty, got %q", result.Input.Arg("command"))
	}
}

func TestParser_BuildSimpleParsed_NotFound(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	p := newTestParser(reg)
	_, err := p.buildSimpleParsed("nonexistent")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestParser_Parse_InvalidFlag(t *testing.T) {
	t.Parallel()
	p, reg := setupParserWithHelp()
	cmd := &mockCommand{name: "cmd"}
	_ = reg.register(cmd)
	_, err := p.parse(&bytes.Buffer{}, []string{"cmd", "--unknown"})
	if err == nil {
		t.Errorf("expected error for unknown flag")
	}
}
