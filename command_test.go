package cli

import (
	"testing"
)

func TestStringArg_Basic(t *testing.T) {
	t.Parallel()
	arg := StringArg("file", "Input file")
	if arg.Name != "file" {
		t.Errorf("expected name %q, got %q", "file", arg.Name)
	}
	if arg.Description != "Input file" {
		t.Errorf("expected description %q, got %q", "Input file", arg.Description)
	}
	if arg.Optional {
		t.Errorf("expected Optional=false, got true")
	}
	if arg.Default != "" {
		t.Errorf("expected empty default, got %q", arg.Default)
	}
}

func TestStringArgOptional_Basic(t *testing.T) {
	t.Parallel()
	arg := StringArgOptional("file", "default.txt", "Input file")
	if arg.Name != "file" {
		t.Errorf("expected name %q, got %q", "file", arg.Name)
	}
	if arg.Description != "Input file" {
		t.Errorf("expected description %q, got %q", "Input file", arg.Description)
	}
	if !arg.Optional {
		t.Errorf("expected Optional=true, got false")
	}
	if arg.Default != "default.txt" {
		t.Errorf("expected default %q, got %q", "default.txt", arg.Default)
	}
}

func TestStringOption_Basic(t *testing.T) {
	t.Parallel()
	opt := StringOption("output", "o", "out.txt", "Output file")
	if opt.Name != "output" {
		t.Errorf("expected name %q, got %q", "output", opt.Name)
	}
	if opt.Short != "o" {
		t.Errorf("expected short %q, got %q", "o", opt.Short)
	}
	if opt.Default != "out.txt" {
		t.Errorf("expected default %q, got %q", "out.txt", opt.Default)
	}
	if opt.Type != OptionTypeString {
		t.Errorf("expected type %d, got %d", OptionTypeString, opt.Type)
	}
}

func TestIntOption_Basic(t *testing.T) {
	t.Parallel()
	opt := IntOption("count", "c", 10, "Count value")
	if opt.Name != "count" {
		t.Errorf("expected name %q, got %q", "count", opt.Name)
	}
	if opt.Short != "c" {
		t.Errorf("expected short %q, got %q", "c", opt.Short)
	}
	if opt.Default != 10 {
		t.Errorf("expected default %d, got %v", 10, opt.Default)
	}
	if opt.Type != OptionTypeInt {
		t.Errorf("expected type %d, got %d", OptionTypeInt, opt.Type)
	}
}

func TestBoolOption_Basic(t *testing.T) {
	t.Parallel()
	opt := BoolOption("verbose", "V", true, "Verbose mode")
	if opt.Name != "verbose" {
		t.Errorf("expected name %q, got %q", "verbose", opt.Name)
	}
	if opt.Short != "V" {
		t.Errorf("expected short %q, got %q", "V", opt.Short)
	}
	if opt.Default != true {
		t.Errorf("expected default true, got %v", opt.Default)
	}
	if opt.Type != OptionTypeBool {
		t.Errorf("expected type %d, got %d", OptionTypeBool, opt.Type)
	}
}
