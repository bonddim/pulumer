package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

// run executes a git command and returns the output
func (g Client) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

// getRoot returns the root directory of the git repository
func (g Client) getRoot() (string, error) {
	output, err := g.run("rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("failed to get repository root directory: %w", err)
	}
	return output, nil
}

// Diff returns a list of files changed since the baseRef
func (g Client) Diff(baseRef string) ([]string, error) {
	output, err := g.run("diff", "--name-only", "--merge-base", baseRef)
	if err != nil {
		// If the merge-base fails, try to get the diff without it
		output, err = g.run("diff", "--name-only", baseRef)
		if err != nil {
			return nil, fmt.Errorf("failed to get diff with '%s': %w", baseRef, err)
		}
	}
	return strings.Split(output, "\n"), nil
}
