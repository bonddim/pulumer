package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bonddim/pulumer/internal/config"
)

func Diff(cfg config.Config) []string {
	// Run the git command to get the list of changed files
	cmd := exec.Command("git", "diff", "--name-only", "--merge-base", cfg.Base)
	cmd.Dir = cfg.WorkDir
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// Split the output into lines (file names)
	files := strings.Split(strings.TrimSpace(string(output)), "\n")

	return files
}
