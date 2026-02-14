package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestVersionCommand_Metadata(t *testing.T) {
	t.Parallel()
	v := &versionCommand{appName: "myapp", appVersion: "2.0"}
	if v.Name() != cmdNameVersion {
		t.Errorf("expected %s, got %s", cmdNameVersion, v.Name())
	}
	if v.Description() == "" {
		t.Errorf("expected non-empty description")
	}
	if v.Group() != groupConsole {
		t.Errorf("expected %s, got %s", groupConsole, v.Group())
	}
	if v.Args() != nil {
		t.Errorf("expected nil args")
	}
	if v.Options() != nil {
		t.Errorf("expected nil options")
	}
}

func TestVersionCommand_Execute_Success(t *testing.T) {
	t.Parallel()
	v := &versionCommand{appName: "myapp", appVersion: "2.0"}
	var buf bytes.Buffer
	err := v.Execute(context.Background(), nil, &buf, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	expected := "myapp version 2.0"
	if !strings.Contains(buf.String(), expected) {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestVersionCommand_Execute_WriteError(t *testing.T) {
	t.Parallel()
	v := &versionCommand{appName: "myapp", appVersion: "2.0"}
	w := &failAfterNWriter{n: 0}
	err := v.Execute(context.Background(), nil, w, nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
