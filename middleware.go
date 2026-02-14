package cli

import (
	"context"
	"io"
)

type Handler func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error

type Middleware func(next Handler) Handler
