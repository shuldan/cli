package cli

import (
	"context"
	"io"
	"testing"
)

func TestMiddleware_TypeAssignment(t *testing.T) {
	t.Parallel()
	var mw Middleware = func(next Handler) Handler {
		return func(ctx context.Context, in io.Reader, out io.Writer, input *Input) error {
			return next(ctx, in, out, input)
		}
	}
	var h Handler = func(_ context.Context, _ io.Reader, _ io.Writer, _ *Input) error {
		return nil
	}
	wrapped := mw(h)
	err := wrapped(context.Background(), nil, nil, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
