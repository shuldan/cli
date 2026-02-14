package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type parsed struct {
	Command Command
	Input   *Input
}

type parser struct {
	registry *registry
}

func (p *parser) parse(output io.Writer, args []string) (*parsed, error) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		return p.buildHelpParsed("")
	}

	if args[0] == "--version" || args[0] == "-v" {
		return p.buildSimpleParsed("version")
	}

	commandName := args[0]
	commandArgs := args[1:]

	command, exists := p.registry.get(commandName)
	if !exists {
		return nil, &CommandNotFoundError{Name: commandName}
	}

	flagSet := flag.NewFlagSet(commandName, flag.ContinueOnError)
	flagSet.SetOutput(output)
	flagSet.Usage = func() {}

	optionValues, err := p.bindOptions(flagSet, command.Options())
	if err != nil {
		return nil, err
	}

	if err := flagSet.Parse(commandArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return p.buildHelpParsed(commandName)
		}
		return nil, err
	}

	argValues, remaining, err := p.bindArgs(commandName, command.Args(), flagSet.Args())
	if err != nil {
		return nil, err
	}

	options := p.collectOptions(command.Options(), optionValues)

	return &parsed{
		Command: command,
		Input:   newInput(argValues, options, remaining),
	}, nil
}

func (p *parser) bindOptions(flagSet *flag.FlagSet, options []Option) (map[string]any, error) {
	pointers := make(map[string]any, len(options))

	for _, opt := range options {
		switch opt.Type {
		case OptionTypeString:
			def, ok := opt.Default.(string)
			if !ok {
				return nil, fmt.Errorf("option %q: default value is not a string", opt.Name)
			}
			ptr := flagSet.String(opt.Name, def, opt.Description)
			if opt.Short != "" {
				flagSet.StringVar(ptr, opt.Short, def, opt.Description)
			}
			pointers[opt.Name] = ptr

		case OptionTypeInt:
			def, ok := opt.Default.(int)
			if !ok {
				return nil, fmt.Errorf("option %q: default value is not an int", opt.Name)
			}
			ptr := flagSet.Int(opt.Name, def, opt.Description)
			if opt.Short != "" {
				flagSet.IntVar(ptr, opt.Short, def, opt.Description)
			}
			pointers[opt.Name] = ptr

		case OptionTypeBool:
			def, ok := opt.Default.(bool)
			if !ok {
				return nil, fmt.Errorf("option %q: default value is not a bool", opt.Name)
			}
			ptr := flagSet.Bool(opt.Name, def, opt.Description)
			if opt.Short != "" {
				flagSet.BoolVar(ptr, opt.Short, def, opt.Description)
			}
			pointers[opt.Name] = ptr

		default:
			return nil, fmt.Errorf("option %q: unsupported option type %d", opt.Name, opt.Type)
		}
	}

	return pointers, nil
}

func (p *parser) bindArgs(commandName string, declared []Arg, positional []string) (map[string]string, []string, error) {
	argValues := make(map[string]string, len(declared))
	matched := 0

	for i, arg := range declared {
		switch {
		case i < len(positional):
			argValues[arg.Name] = positional[i]
			matched++
		case arg.Optional:
			argValues[arg.Name] = arg.Default
		default:
			return nil, nil, &MissingArgumentError{
				Command:  commandName,
				Argument: arg.Name,
			}
		}
	}

	var remaining []string
	if matched < len(positional) {
		remaining = positional[matched:]
	}

	return argValues, remaining, nil
}

func (p *parser) collectOptions(options []Option, pointers map[string]any) map[string]any {
	result := make(map[string]any, len(options))

	for _, opt := range options {
		ptr := pointers[opt.Name]

		switch opt.Type {
		case OptionTypeString:
			if v, ok := ptr.(*string); ok {
				result[opt.Name] = *v
			}
		case OptionTypeInt:
			if v, ok := ptr.(*int); ok {
				result[opt.Name] = *v
			}
		case OptionTypeBool:
			if v, ok := ptr.(*bool); ok {
				result[opt.Name] = *v
			}
		}
	}

	return result
}

func (p *parser) buildHelpParsed(commandName string) (*parsed, error) {
	helpCmd, exists := p.registry.get("help")
	if !exists {
		return nil, &CommandNotFoundError{Name: "help"}
	}

	args := make(map[string]string)
	if commandName != "" {
		args["command"] = commandName
	}

	return &parsed{
		Command: helpCmd,
		Input:   newInput(args, make(map[string]any), nil),
	}, nil
}

func (p *parser) buildSimpleParsed(commandName string) (*parsed, error) {
	cmd, exists := p.registry.get(commandName)
	if !exists {
		return nil, &CommandNotFoundError{Name: commandName}
	}

	return &parsed{
		Command: cmd,
		Input:   newInput(make(map[string]string), make(map[string]any), nil),
	}, nil
}
