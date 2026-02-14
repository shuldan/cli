package cli

import (
	"testing"
)

func TestInput_Arg_Exists(t *testing.T) {
	t.Parallel()
	inp := newInput(map[string]string{"name": "val"}, nil, nil)
	if got := inp.Arg("name"); got != "val" {
		t.Errorf("expected %q, got %q", "val", got)
	}
}

func TestInput_Arg_NotExists(t *testing.T) {
	t.Parallel()
	inp := newInput(map[string]string{}, nil, nil)
	if got := inp.Arg("missing"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestInput_StringOption_Exists(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{"key": "val"}, nil)
	if got := inp.StringOption("key"); got != "val" {
		t.Errorf("expected %q, got %q", "val", got)
	}
}

func TestInput_StringOption_WrongType(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{"key": 123}, nil)
	if got := inp.StringOption("key"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestInput_StringOption_NotExists(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{}, nil)
	if got := inp.StringOption("missing"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestInput_IntOption_Exists(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{"n": 42}, nil)
	if got := inp.IntOption("n"); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestInput_IntOption_WrongType(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{"n": "str"}, nil)
	if got := inp.IntOption("n"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestInput_IntOption_NotExists(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{}, nil)
	if got := inp.IntOption("missing"); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestInput_BoolOption_Exists(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{"flag": true}, nil)
	if got := inp.BoolOption("flag"); !got {
		t.Errorf("expected true, got false")
	}
}

func TestInput_BoolOption_WrongType(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{"flag": "yes"}, nil)
	if got := inp.BoolOption("flag"); got {
		t.Errorf("expected false, got true")
	}
}

func TestInput_BoolOption_NotExists(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, map[string]any{}, nil)
	if got := inp.BoolOption("missing"); got {
		t.Errorf("expected false, got true")
	}
}

func TestInput_RemainingArgs(t *testing.T) {
	t.Parallel()
	rem := []string{"a", "b"}
	inp := newInput(nil, nil, rem)
	got := inp.RemainingArgs()
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("expected %v, got %v", rem, got)
	}
}

func TestInput_RemainingArgs_Nil(t *testing.T) {
	t.Parallel()
	inp := newInput(nil, nil, nil)
	got := inp.RemainingArgs()
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestNewInput_AllFields(t *testing.T) {
	t.Parallel()
	args := map[string]string{"a": "1"}
	opts := map[string]any{"b": 2}
	rem := []string{"c"}
	inp := newInput(args, opts, rem)
	if inp.Arg("a") != "1" {
		t.Errorf("expected arg a=1")
	}
	if inp.IntOption("b") != 2 {
		t.Errorf("expected opt b=2")
	}
	if len(inp.RemainingArgs()) != 1 {
		t.Errorf("expected 1 remaining")
	}
}
