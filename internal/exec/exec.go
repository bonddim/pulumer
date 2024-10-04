package exec

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type Exec struct {
	ctx context.Context
	cmd auto.PulumiCommand
}

func NewExec() *Exec {
	ctx := context.Background()
	cmd := getPulumi(ctx)

	return &Exec{cmd: cmd, ctx: ctx}
}

func getPulumi(ctx context.Context) auto.PulumiCommand {
	cmd, err := auto.NewPulumiCommand(&auto.PulumiCommandOptions{})
	if err != nil {
		cmd, err = auto.InstallPulumiCommand(ctx, &auto.PulumiCommandOptions{})
		if err != nil {
			fmt.Printf("unable to initialize pulumi cli\n%v", err)
			os.Exit(1)
		}
	}
	return cmd
}

// Run runs a command with the given arguments, attaches stdout/stderr, and exits the main program on failure.
func (exec Exec) Run(workdir string, args ...string) {
	var stdin io.Reader = os.Stdin
	var stdoutWriters []io.Writer = []io.Writer{os.Stdout}
	var stderrWriters []io.Writer = []io.Writer{os.Stderr}
	var envVars []string = os.Environ()

	_, _, code, err := exec.cmd.Run(exec.ctx, workdir, stdin, stdoutWriters, stderrWriters, envVars, args...)
	if err != nil {
		fmt.Printf("Failed to run command %v\n", err)
		os.Exit(code)
	}
}
