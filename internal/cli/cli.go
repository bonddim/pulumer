package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/bonddim/pulumer/internal/config"
	"github.com/bonddim/pulumer/internal/git"
	"github.com/bonddim/pulumer/internal/pulumi"
	"github.com/bonddim/pulumer/internal/stacks"
)

type CLI struct {
	cfg       config.Config
	ctx       context.Context
	cwd       string
	git       *git.Client
	pulumiCli *pulumi.CLI
	stacks    stacks.Stacks
}

func NewCLI(ctx context.Context, cfgFile string) (*CLI, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		return nil, err
	}

	cli := &CLI{
		cfg:    *cfg,
		ctx:    ctx,
		cwd:    cwd,
		git:    git.NewClient(),
		stacks: stacks.NewStacks(),
	}

	return cli, nil
}

func (c *CLI) getStacks() error {
	// search for stack files
	stackFiles, err := c.git.SearchFiles(c.cwd, `^Pulumi\..*\.yaml$`)
	if err != nil {
		return err
	}

	// get stack from the working directory
	if err := c.stacks.New(stackFiles); err != nil {
		return fmt.Errorf("failed to get stacks: %w", err)
	}

	// sort stacks
	if c.cfg.Sorted {
		if err := c.stacks.Sort(); err != nil {
			return fmt.Errorf("failed to sort stacks: %w", err)
		}
	}

	// swap stacks
	if c.cfg.Reversed {
		c.stacks.Reverse()
	}

	// filter changed stacks
	if c.cfg.Changed {
		diff, err := c.git.Diff(c.cfg.Since)
		if err != nil {
			return err
		}
		c.stacks.Changed(diff)
	}

	return nil
}
