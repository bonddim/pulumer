package cli

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/bonddim/pulumer/internal/pulumi"
)

// Run runs pulumi command with the given arguments, attaches stdout/stderr, and exits the main program on failure.
func (c *CLI) Run(args []string) error {
	// init pulumi cli
	pulumiCli, err := pulumi.NewCLI(c.ctx, c.cfg.PulumiVersion)
	if err != nil {
		return err
	}
	c.pulumiCli = pulumiCli

	stackCommands := []string{
		"cancel",
		"destroy",
		"preview",
		"refresh",
		"up",
	}

	if slices.Contains(stackCommands, args[0]) {
		return c.processStacks(args)
	}

	return c.pulumiCli.Run(c.cwd, args)
}

// processStacks runs pulumi stack commands.
func (c *CLI) processStacks(args []string) error {
	// Run stack commands on sorted stacks by default.
	c.cfg.Sorted = true
	// Reverse the order for destroy.
	if args[0] == "destroy" {
		c.cfg.Reversed = true
	}

	if err := c.getStacks(); err != nil {
		return err
	}

	pulumiArgs := getArgs(args)
	workspaceInstalledDeps := []string{}
	for _, stack := range c.stacks {
		fmt.Fprintf(os.Stderr, "Running: pulumi %s for stack %s ...\n", strings.Join(args, " "), stack.Id)
		// Install dependencies for each workspace if not installed already
		if !slices.Contains(workspaceInstalledDeps, stack.Workspace) {
			if err := c.pulumiCli.Run(stack.Workspace, []string{"install", "--use-language-version-tools"}); err != nil {
				return err
			}
			workspaceInstalledDeps = append(workspaceInstalledDeps, stack.Workspace)
		}

		if err := c.pulumiCli.Run(stack.Workspace, append(pulumiArgs, "--stack", stack.Name)); err != nil {
			return err
		}
	}

	return nil
}

func getArgs(args []string) []string {
	const skipPreview = "--skip-preview"
	const suppressProgress = "--suppress-progress"

	switch args[0] {
	case "preview":
		return append(args, suppressProgress, "--diff")
	case "refresh":
		return append(args, suppressProgress, skipPreview)
	case "up":
		return append(args, suppressProgress, skipPreview)
	case "destroy":
		return append(args, suppressProgress, skipPreview)
	default:
		return args
	}
}
