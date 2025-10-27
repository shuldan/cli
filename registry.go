package cli

import (
	"errors"
	"sort"
	"sync"
)

type registry struct {
	mu       sync.RWMutex
	commands map[string]Command
	groups   map[string][]string
}

func (r *registry) register(command Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if command == nil {
		return errors.New("command is nil")
	}

	name := command.Name()
	if name == "" {
		return errors.New("command name is empty")
	}

	if _, exists := r.commands[name]; exists {
		return errors.New("command exists")
	}

	r.commands[name] = command

	group := command.Group()
	if group == "" {
		group = "general"
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

func (r *registry) getGroups() map[string][]Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string][]Command)
	for group, names := range r.groups {
		commands := make([]Command, 0, len(names))
		for _, name := range names {
			if command, exists := r.commands[name]; exists && command != nil {
				commands = append(commands, command)
			}
		}

		sort.Slice(commands, func(i, j int) bool {
			return commands[i].Name() < commands[j].Name()
		})

		result[group] = commands
	}
	return result
}
