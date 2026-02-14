package cli

import (
	"context"
	"io"
	"runtime/debug"
)

type executor struct {
	parser     *parser
	middleware []Middleware
}

func (e *executor) execute(ctx context.Context, in io.Reader, out io.Writer, args []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return e.withRecovery(func() error {
		p, err := e.parser.parse(out, args)
		if err != nil {
			return err
		}

		handler := func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error {
			return p.Command.Execute(ctx, in, out, input)
		}

		for i := len(e.middleware) - 1; i >= 0; i-- {
			handler = e.middleware[i](handler)
		}

		return handler(ctx, in, out, p.Input)
	})
}

func (e *executor) withRecovery(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{
				Value: r,
				Stack: debug.Stack(),
			}
		}
	}()
	return fn()
}
