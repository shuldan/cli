package cli

import (
	"context"
	"io"
)

type Command interface {
	Name() string
	Description() string
	Group() string
	Args() []Arg
	Options() []Option
	Execute(ctx context.Context, in io.Reader, out io.Writer, input *Input) error
}

type Arg struct {
	Name        string
	Description string
	Optional    bool
	Default     string
}

func StringArg(name, description string) Arg {
	return Arg{
		Name:        name,
		Description: description,
	}
}

func StringArgOptional(name, defaultVal, description string) Arg {
	return Arg{
		Name:        name,
		Description: description,
		Optional:    true,
		Default:     defaultVal,
	}
}

type OptionType int

const (
	OptionTypeString OptionType = iota
	OptionTypeInt
	OptionTypeBool
)

type Option struct {
	Name        string
	Short       string
	Description string
	Default     any
	Type        OptionType
}

func StringOption(name, short, defaultVal, description string) Option {
	return Option{
		Name:        name,
		Short:       short,
		Description: description,
		Default:     defaultVal,
		Type:        OptionTypeString,
	}
}

func IntOption(name, short string, defaultVal int, description string) Option {
	return Option{
		Name:        name,
		Short:       short,
		Description: description,
		Default:     defaultVal,
		Type:        OptionTypeInt,
	}
}

func BoolOption(name, short string, defaultVal bool, description string) Option {
	return Option{
		Name:        name,
		Short:       short,
		Description: description,
		Default:     defaultVal,
		Type:        OptionTypeBool,
	}
}
