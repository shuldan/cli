package cli

import (
	"context"
	"fmt"
	"io"
)

type helpCommand struct {
	registry *registry
	appName  string
}

func (h *helpCommand) Name() string        { return "help" }
func (h *helpCommand) Description() string { return "Display help for commands" }
func (h *helpCommand) Group() string       { return "console" }

func (h *helpCommand) Args() []Arg {
	return []Arg{
		StringArgOptional("command", "", "Command to show help for"),
	}
}

func (h *helpCommand) Options() []Option { return nil }

func (h *helpCommand) Execute(_ context.Context, _ io.Reader, out io.Writer, input *Input) error {
	cmdName := input.Arg("command")
	if cmdName != "" {
		return h.showCommandHelp(cmdName, out)
	}
	return h.showGeneralHelp(out)
}

func (h *helpCommand) showGeneralHelp(out io.Writer) error {
	groups := h.registry.getGroups()

	if _, err := fmt.Fprintf(out, "%s\n\nUsage: %s <command> [options] [arguments]\n\n", h.appName, h.appName); err != nil {
		return err
	}

	for _, group := range groups {
		maxLen := 0
		for _, cmd := range group.Commands {
			if len(cmd.Name()) > maxLen {
				maxLen = len(cmd.Name())
			}
		}

		if _, err := fmt.Fprintf(out, "%s:\n", group.Name); err != nil {
			return err
		}

		for _, cmd := range group.Commands {
			if _, err := fmt.Fprintf(out, "  %-*s    %s\n", maxLen, cmd.Name(), cmd.Description()); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(out); err != nil {
			return err
		}
	}

	return nil
}

func (h *helpCommand) showCommandHelp(cmdName string, out io.Writer) error {
	cmd, exists := h.registry.get(cmdName)
	if !exists {
		return &CommandNotFoundError{Name: cmdName}
	}

	usage := h.buildUsageLine(cmd)

	if _, err := fmt.Fprintf(out, "%s — %s\n\nUsage: %s\n", cmd.Name(), cmd.Description(), usage); err != nil {
		return err
	}

	if err := h.printArgs(out, cmd.Args()); err != nil {
		return err
	}

	if err := h.printOptions(out, cmd.Options()); err != nil {
		return err
	}

	return nil
}

func (h *helpCommand) buildUsageLine(cmd Command) string {
	usage := h.appName + " " + cmd.Name()

	for _, arg := range cmd.Args() {
		if arg.Optional {
			usage += " [" + arg.Name + "]"
		} else {
			usage += " <" + arg.Name + ">"
		}
	}

	for _, opt := range cmd.Options() {
		if opt.Type == OptionTypeBool {
			usage += " [--" + opt.Name + "]"
		} else {
			usage += " [--" + opt.Name + "=...]"
		}
	}

	return usage
}

func (h *helpCommand) printArgs(out io.Writer, args []Arg) error {
	if len(args) == 0 {
		return nil
	}

	if _, err := fmt.Fprintf(out, "\nArguments:\n"); err != nil {
		return err
	}

	maxLen := 0
	for _, a := range args {
		if len(a.Name) > maxLen {
			maxLen = len(a.Name)
		}
	}

	for _, a := range args {
		desc := a.Description
		if a.Optional && a.Default != "" {
			desc += fmt.Sprintf(" (default: %s)", a.Default)
		}
		if _, err := fmt.Fprintf(out, "  %-*s    %s\n", maxLen, a.Name, desc); err != nil {
			return err
		}
	}

	return nil
}

func (h *helpCommand) printOptions(out io.Writer, options []Option) error {
	if len(options) == 0 {
		return nil
	}

	if _, err := fmt.Fprintf(out, "\nOptions:\n"); err != nil {
		return err
	}

	type entry struct {
		label string
		desc  string
	}

	entries := make([]entry, 0, len(options))
	maxLen := 0

	for _, opt := range options {
		label := "--" + opt.Name
		if opt.Short != "" {
			label += ", -" + opt.Short
		}
		if len(label) > maxLen {
			maxLen = len(label)
		}

		desc := opt.Description
		if opt.Default != nil {
			desc += fmt.Sprintf(" (default: %v)", opt.Default)
		}

		entries = append(entries, entry{label: label, desc: desc})
	}

	for _, e := range entries {
		if _, err := fmt.Fprintf(out, "  %-*s    %s\n", maxLen, e.label, e.desc); err != nil {
			return err
		}
	}

	return nil
}
