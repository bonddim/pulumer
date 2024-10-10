/*
Copyright © 2024 Dmytro Bondar <git@bonddim.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/bonddim/pulumer/internal/exec"
	"github.com/bonddim/pulumer/internal/git"
	"github.com/bonddim/pulumer/internal/stacks"
	"github.com/spf13/cobra"
)

var (
	cli   *exec.Exec
	login bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pulumi command in the stacks",
	Args:  cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cli = exec.NewExec()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if login {
			if os.Getenv("PULUMI_BACKEND_URL") == "" && os.Getenv("PULUMI_ACCESS_TOKEN") == "" {
				cobra.CheckErr("PULUMI_BACKEND_URL or PULUMI_ACCESS_TOKEN environment variable is required to login.")
			}
			cli.Run(cfg.WorkDir, "login")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		commands := []string{
			"up",
			"destroy",
			"preview",
			"refresh",
		}

		if slices.Contains(commands, args[0]) {
			changed, _ := cmd.Flags().GetBool("changed")

			stacks, err := stacks.GetStacks(cfg)
			cobra.CheckErr(err)

			err = stacks.Sorted()
			cobra.CheckErr(err)

			if changed {
				gitDiff := git.Diff(cfg)
				stacks = stacks.FilterChanged(gitDiff)
			}

			if args[0] == "destroy" {
				stacks = stacks.Reverse()
			}

			for _, stack := range stacks {
				stackArgs := append(args, "--stack", stack.Name)
				if args[0] != "preview" {
					stackArgs = append(stackArgs, "--skip-preview")
				}
				cli.Run(filepath.Join(cfg.WorkDir, stack.Workspace), stackArgs...)
			}

			return
		}

		cli.Run(cfg.WorkDir, args...)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("changed", "c", false, "Filter stacks based on changes made in git")
	runCmd.Flags().BoolVarP(&login, "login", "l", false, "Login to pulumi backend. Requires PULUMI_BACKEND_URL or PULUMI_ACCESS_TOKEN environment variable")
}
