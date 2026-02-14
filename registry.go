package cli

import (
	"fmt"
	"sort"
	"sync"
)

type commandGroup struct {
	Name     string
	Commands []Command
}

type registry struct {
	mu       sync.RWMutex
	commands map[string]Command
	groups   map[string][]string
	order    []string
}

func (r *registry) register(command Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if command == nil {
		return &InvalidCommandError{Command: "", Reason: "command is nil"}
	}

	name := command.Name()
	if name == "" {
		return &InvalidCommandError{Command: "", Reason: "command name is empty"}
	}

	if _, exists := r.commands[name]; exists {
		return &CommandExistsError{Name: name}
	}

	if err := r.validateArgs(command); err != nil {
		return err
	}

	if err := r.validateOptions(command); err != nil {
		return err
	}

	r.commands[name] = command

	group := command.Group()
	if group == "" {
		group = groupGeneral
	}

	if _, exists := r.groups[group]; !exists {
		r.order = append(r.order, group)
	}

	r.groups[group] = append(r.groups[group], name)

	return nil
}

func (r *registry) get(name string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	command, exists := r.commands[name]
	return command, exists
}

func (r *registry) getGroups() []commandGroup {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]commandGroup, 0, len(r.order))

	for _, groupName := range r.order {
		names := r.groups[groupName]
		commands := make([]Command, 0, len(names))

		for _, name := range names {
			if cmd, exists := r.commands[name]; exists && cmd != nil {
				commands = append(commands, cmd)
			}
		}

		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Name() < commands[j].Name()
		})

		result = append(result, commandGroup{
			Name:     groupName,
			Commands: commands,
		})
	}

	return result
}

func (r *registry) validateArgs(command Command) error {
	args := command.Args()
	if args == nil {
		return nil
	}

	seenOptional := false
	seen := make(map[string]bool)

	for _, arg := range args {
		if arg.Name == "" {
			return &InvalidCommandError{
				Command: command.Name(),
				Reason:  "argument name is empty",
			}
		}

		if seen[arg.Name] {
			return &InvalidCommandError{
				Command: command.Name(),
				Reason:  fmt.Sprintf("duplicate argument name: %s", arg.Name),
			}
		}
		seen[arg.Name] = true

		if arg.Optional {
			seenOptional = true
		} else if seenOptional {
			return &InvalidCommandError{
				Command: command.Name(),
				Reason:  fmt.Sprintf("required argument %q cannot follow optional argument", arg.Name),
			}
		}
	}

	return nil
}

func (r *registry) validateOptions(command Command) error {
	options := command.Options()
	if options == nil {
		return nil
	}

	names := make(map[string]bool)
	shorts := make(map[string]bool)

	for _, opt := range options {
		if opt.Name == "" {
			return &InvalidCommandError{
				Command: command.Name(),
				Reason:  "option name is empty",
			}
		}

		if names[opt.Name] {
			return &InvalidCommandError{
				Command: command.Name(),
				Reason:  fmt.Sprintf("duplicate option name: %s", opt.Name),
			}
		}
		names[opt.Name] = true

		if opt.Short != "" {
			if shorts[opt.Short] {
				return &InvalidCommandError{
					Command: command.Name(),
					Reason:  fmt.Sprintf("duplicate short option: %s", opt.Short),
				}
			}
			shorts[opt.Short] = true
		}
	}

	return nil
}
