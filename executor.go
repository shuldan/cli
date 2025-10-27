package cli

import (
	"context"
	"fmt"
	"io"
)

type executor struct {
	parser *parser
}

func (e *executor) execute(ctx context.Context, input io.Reader, output io.Writer, args []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := e.withRecovery(func() error {
		parsed, err := e.parser.parse(output, args)
		if err != nil {
			return err
		}
		return parsed.Command.Execute(ctx, input, output, args)
	}); err != nil {
		return err
	}

	return nil
}

func (e *executor) withRecovery(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during command execution: %v", r)
		}
	}()
	return fn()
}
