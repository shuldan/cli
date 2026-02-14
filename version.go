package cli

import (
	"context"
	"fmt"
	"io"
)

type versionCommand struct {
	appName    string
	appVersion string
}

func (v *versionCommand) Name() string        { return "version" }
func (v *versionCommand) Description() string { return "Display application version" }
func (v *versionCommand) Group() string       { return "console" }
func (v *versionCommand) Args() []Arg         { return nil }
func (v *versionCommand) Options() []Option   { return nil }

func (v *versionCommand) Execute(_ context.Context, _ io.Reader, out io.Writer, _ *Input) error {
	_, err := fmt.Fprintf(out, "%s version %s\n", v.appName, v.appVersion)
	return err
}
