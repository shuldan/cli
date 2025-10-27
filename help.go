package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"text/template"
)

type help struct {
	registry *registry
	command  string
}

func (h *help) Name() string {
	return "help"
}

func (h *help) Description() string {
	return "Display help for commands"
}

func (h *help) Group() string {
	return "console"
}

func (h *help) Configure(flags *flag.FlagSet) {
	flags.StringVar(&h.command, "command", "", "Show help for specific command")
}

func (h *help) Execute(ctx context.Context, _ io.Reader, output io.Writer, _ []string) error {
	if h.command != "" {
		if err := h.showCommandHelp(h.command, output); err != nil {
			return err
		}
		ctx.Done()
		return nil
	}

	if err := h.showGeneralHelp(output); err != nil {
		return err
	}
	ctx.Done()
	return nil
}

func (h *help) showGeneralHelp(output io.Writer) error {
	groups := h.registry.getGroups()

	data := struct {
		Groups map[string][]PrintableCommand
	}{
		Groups: make(map[string][]PrintableCommand),
	}

	for groupName, commands := range groups {
		longest := 0
		for _, cmd := range commands {
			if len(cmd.Name()) > longest {
				longest = len(cmd.Name())
			}
		}

		formatter := "%-" + strconv.Itoa(longest) + "s"
		printableCommands := make([]PrintableCommand, 0, len(commands))

		for _, cmd := range commands {
			printableCommands = append(printableCommands, PrintableCommand{
				PaddedName:  fmt.Sprintf(formatter, cmd.Name()),
				Description: cmd.Description(),
			})
		}

		data.Groups[groupName] = printableCommands
	}

	tmpl := template.New("help")
	helpTemplate := `Usage: command [options] [arguments]

{{ range $group, $commands := .Groups }}{{ $group }}:{{ range $commands }}
  {{.PaddedName}}  {{.Description}}{{ end }}

{{ end }}`

	template.Must(tmpl.Parse(helpTemplate))

	return tmpl.Execute(output, data)
}

func (h *help) showCommandHelp(commandName string, output io.Writer) error {
	command, exists := h.registry.get(commandName)
	if !exists || command == nil {
		return errors.New("command not found")
	}

	if _, err := fmt.Fprintf(output, "%s - %s\n\n", command.Name(), command.Description()); err != nil {
		return err
	}

	flags := flag.NewFlagSet(command.Name(), flag.ContinueOnError)
	flags.SetOutput(output)
	command.Configure(flags)

	if _, err := fmt.Fprintf(output, "Options:\n"); err != nil {
		return err
	}

	flags.PrintDefaults()

	return nil
}

type PrintableCommand struct {
	PaddedName  string
	Description string
}
