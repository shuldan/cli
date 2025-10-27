package cli

import (
	"errors"
	"flag"
	"io"
)

type parsedCommand struct {
	Name    string
	Args    []string
	Flags   *flag.FlagSet
	Command Command
}

type parser struct {
	registry *registry
}

func (p *parser) parse(output io.Writer, args []string) (*parsedCommand, error) {
	if len(args) == 0 {
		return nil, errors.New("no command specified")
	}

	commandName := args[0]
	commandArgs := args[1:]

	command, exists := p.registry.get(commandName)
	if !exists {
		return nil, errors.New("unknown command: " + commandName)
	}

	flagSet := flag.NewFlagSet(commandName, flag.ContinueOnError)
	flagSet.SetOutput(output)

	command.Configure(flagSet)

	if err := flagSet.Parse(commandArgs); err != nil {
		return nil, err
	}

	return &parsedCommand{
		Name:    commandName,
		Args:    flagSet.Args(),
		Flags:   flagSet,
		Command: command,
	}, nil
}
