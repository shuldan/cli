package cli

type Input struct {
	args          map[string]string
	options       map[string]any
	remainingArgs []string
}

func newInput(args map[string]string, options map[string]any, remaining []string) *Input {
	return &Input{
		args:          args,
		options:       options,
		remainingArgs: remaining,
	}
}

func (i *Input) Arg(name string) string {
	if v, ok := i.args[name]; ok {
		return v
	}
	return ""
}

func (i *Input) StringOption(name string) string {
	if v, ok := i.options[name]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (i *Input) IntOption(name string) int {
	if v, ok := i.options[name]; ok {
		if n, ok := v.(int); ok {
			return n
		}
	}
	return 0
}

func (i *Input) BoolOption(name string) bool {
	if v, ok := i.options[name]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

func (i *Input) RemainingArgs() []string {
	return i.remainingArgs
}
