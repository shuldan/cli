package cli

import (
	"context"
	"testing"
)

func TestSignalContext_Basic(t *testing.T) {
	t.Parallel()
	ctx, cancel := SignalContext(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		t.Errorf("expected context to not be done yet")
	default:
	}
	cancel()
	<-ctx.Done()
}
