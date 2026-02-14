package cli

import (
	"errors"
	"fmt"
	"testing"
)

func TestGetExitCode_Table(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"nil error", nil, ExitSuccess},
		{"ExitCoder", &ExitError{Code: 42}, 42},
		{"CommandNotFound", &CommandNotFoundError{Name: "x"}, ExitUsageError},
		{"MissingArgument", &MissingArgumentError{Command: "c", Argument: "a"}, ExitUsageError},
		{"PanicError", &PanicError{Value: "boom"}, ExitFailure},
		{"generic error", fmt.Errorf("fail"), ExitFailure},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := GetExitCode(tc.err)
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestExitError_Error_WithErr(t *testing.T) {
	t.Parallel()
	e := &ExitError{Code: 1, Err: fmt.Errorf("inner")}
	if e.Error() != "inner" {
		t.Errorf("expected %q, got %q", "inner", e.Error())
	}
}

func TestExitError_Error_NilErr(t *testing.T) {
	t.Parallel()
	e := &ExitError{Code: 5, Err: nil}
	expected := "exit code 5"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestExitError_Unwrap(t *testing.T) {
	t.Parallel()
	inner := fmt.Errorf("root cause")
	e := &ExitError{Code: 1, Err: inner}
	if !errors.Is(e, inner) {
		t.Errorf("expected Unwrap to return inner error")
	}
}

func TestExitError_ExitCode(t *testing.T) {
	t.Parallel()
	e := &ExitError{Code: 77}
	if e.ExitCode() != 77 {
		t.Errorf("expected 77, got %d", e.ExitCode())
	}
}

func TestCommandNotFoundError_Error(t *testing.T) {
	t.Parallel()
	e := &CommandNotFoundError{Name: "foo"}
	expected := "command not found: foo"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestCommandExistsError_Error(t *testing.T) {
	t.Parallel()
	e := &CommandExistsError{Name: "bar"}
	expected := "command already registered: bar"
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestMissingArgumentError_Error(t *testing.T) {
	t.Parallel()
	e := &MissingArgumentError{Command: "cmd", Argument: "arg"}
	expected := `missing required argument <arg> for command "cmd"`
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestInvalidCommandError_Error(t *testing.T) {
	t.Parallel()
	e := &InvalidCommandError{Command: "cmd", Reason: "bad"}
	expected := `invalid command "cmd": bad`
	if e.Error() != expected {
		t.Errorf("expected %q, got %q", expected, e.Error())
	}
}

func TestPanicError_Error(t *testing.T) {
	t.Parallel()
	e := &PanicError{Value: "boom", Stack: []byte("trace")}
	got := e.Error()
	if got == "" {
		t.Errorf("expected non-empty error string")
	}
}
