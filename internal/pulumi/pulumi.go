package pulumi

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/blang/semver"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type CLI struct {
	ctx context.Context
	cmd auto.PulumiCommand
}

func NewCLI(ctx context.Context, version string) (*CLI, error) {
	requiredVersion, err := semver.Parse(version)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}

	cmd, err := getPulumi(ctx, requiredVersion)
	if err != nil {
		return nil, err
	}
	return &CLI{cmd: cmd, ctx: ctx}, nil
}

func getPulumi(ctx context.Context, version semver.Version) (auto.PulumiCommand, error) {
	cmd, err := auto.NewPulumiCommand(&auto.PulumiCommandOptions{
		Version: version,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "pulumi cli with required version not found in PATH, installing...")
		cmd, err = auto.InstallPulumiCommand(ctx, &auto.PulumiCommandOptions{})
		if err != nil {
			return nil, fmt.Errorf("unable to install pulumi cli: %w", err)
		}
	}
	fmt.Fprintf(os.Stderr, "pulumi cli initialized, version: %s\n", cmd.Version().String())
	return cmd, nil
}

// Run runs a command with the given arguments, attaches stdout/stderr, and exits the main program on failure.
func (p CLI) Run(workdir string, args []string) error {
	var stdin io.Reader = os.Stdin
	var stdoutWriters []io.Writer = []io.Writer{os.Stdout}
	var stderrWriters []io.Writer = []io.Writer{os.Stderr}
	var envVars []string = os.Environ()

	_, _, _, err := p.cmd.Run(p.ctx, workdir, stdin, stdoutWriters, stderrWriters, envVars, args...)
	if err != nil {
		return fmt.Errorf("failed pulumi run: %w", err)
	}
	return nil
}
