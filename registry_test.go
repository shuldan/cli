package cli

import (
	"strings"
	"testing"
)

func TestRegistry_Register_Success(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "test", description: "desc"}
	err := reg.register(cmd)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRegistry_Register_NilCommand(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	err := reg.register(nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_Register_EmptyName(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: ""}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_Register_Duplicate(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "dup"}
	_ = reg.register(cmd)
	err := reg.register(&mockCommand{name: "dup"})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_Register_EmptyGroup(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "test", group: ""}
	_ = reg.register(cmd)
	groups := reg.getGroups()
	found := false
	for _, g := range groups {
		if g.Name == "general" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected group 'general'")
	}
}

func TestRegistry_Register_CustomGroup(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "test", group: "mygroup"}
	_ = reg.register(cmd)
	groups := reg.getGroups()
	found := false
	for _, g := range groups {
		if g.Name == "mygroup" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected group 'mygroup'")
	}
}

func TestRegistry_Get_Exists(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "test"}
	_ = reg.register(cmd)
	got, ok := reg.get("test")
	if !ok || got == nil {
		t.Errorf("expected command, got nil")
	}
}

func TestRegistry_Get_NotExists(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	_, ok := reg.get("missing")
	if ok {
		t.Errorf("expected false, got true")
	}
}

func TestRegistry_GetGroups_Sorted(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	_ = reg.register(&mockCommand{name: "beta", group: "g"})
	_ = reg.register(&mockCommand{name: "alpha", group: "g"})
	groups := reg.getGroups()
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Commands[0].Name() != "alpha" {
		t.Errorf("expected alpha first, got %s", groups[0].Commands[0].Name())
	}
}

func TestRegistry_ValidateArgs_EmptyArgName(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "cmd", args: []Arg{{Name: ""}}}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_ValidateArgs_DuplicateArgName(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{{Name: "a"}, {Name: "a"}},
	}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_ValidateArgs_RequiredAfterOptional(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{
			{Name: "opt", Optional: true},
			{Name: "req", Optional: false},
		},
	}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "cannot follow optional") {
		t.Errorf("expected 'cannot follow optional', got %q", err.Error())
	}
}

func TestRegistry_ValidateArgs_NilArgs(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "cmd", args: nil}
	err := reg.register(cmd)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRegistry_ValidateOptions_EmptyOptName(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{
		name:    "cmd",
		options: []Option{{Name: "", Type: OptionTypeString, Default: ""}},
	}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_ValidateOptions_DuplicateOptName(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{
		name: "cmd",
		options: []Option{
			{Name: "a", Type: OptionTypeString, Default: ""},
			{Name: "a", Type: OptionTypeString, Default: ""},
		},
	}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_ValidateOptions_DuplicateShort(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{
		name: "cmd",
		options: []Option{
			{Name: "a", Short: "x", Type: OptionTypeString, Default: ""},
			{Name: "b", Short: "x", Type: OptionTypeString, Default: ""},
		},
	}
	err := reg.register(cmd)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestRegistry_ValidateOptions_NilOptions(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{name: "cmd", options: nil}
	err := reg.register(cmd)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRegistry_ValidateArgs_ValidOptionalOrder(t *testing.T) {
	t.Parallel()
	reg := newTestRegistry()
	cmd := &mockCommand{
		name: "cmd",
		args: []Arg{
			{Name: "req", Optional: false},
			{Name: "opt", Optional: true},
		},
	}
	err := reg.register(cmd)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
