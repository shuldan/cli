package console

import (
	"context"
	"flag"
	"io"
)

type Command interface {
	Name() string
	Description() string
	Group() string
	Configure(flags *flag.FlagSet)
	Execute(ctx context.Context, input io.Reader, output io.Writer, args []string) error
}
