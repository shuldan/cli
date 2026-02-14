package cli

import "fmt"

const (
	ExitSuccess    = 0
	ExitFailure    = 1
	ExitUsageError = 2
)

type ExitCoder interface {
	ExitCode() int
}

func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	if coder, ok := err.(ExitCoder); ok {
		return coder.ExitCode()
	}
	return ExitFailure
}

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("exit code %d", e.Code)
}

func (e *ExitError) Unwrap() error { return e.Err }
func (e *ExitError) ExitCode() int { return e.Code }

type CommandNotFoundError struct {
	Name string
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("command not found: %s", e.Name)
}

func (e *CommandNotFoundError) ExitCode() int { return ExitUsageError }

type CommandExistsError struct {
	Name string
}

func (e *CommandExistsError) Error() string {
	return fmt.Sprintf("command already registered: %s", e.Name)
}

type MissingArgumentError struct {
	Command  string
	Argument string
}

func (e *MissingArgumentError) Error() string {
	return fmt.Sprintf("missing required argument <%s> for command %q", e.Argument, e.Command)
}

func (e *MissingArgumentError) ExitCode() int { return ExitUsageError }

type InvalidCommandError struct {
	Command string
	Reason  string
}

func (e *InvalidCommandError) Error() string {
	return fmt.Sprintf("invalid command %q: %s", e.Command, e.Reason)
}

type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic during command execution: %v\n\n%s", e.Value, e.Stack)
}

func (e *PanicError) ExitCode() int { return ExitFailure }
