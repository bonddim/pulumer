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
	"fmt"

	"github.com/bonddim/pulumer/internal/git"
	"github.com/bonddim/pulumer/internal/stacks"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pulumi stacks",
	Run: func(cmd *cobra.Command, args []string) {
		changed, _ := cmd.Flags().GetBool("changed")
		sorted, _ := cmd.Flags().GetBool("sorted")

		stacks, err := stacks.GetStacks(cfg)
		cobra.CheckErr(err)

		if sorted {
			err = stacks.Sorted()
			cobra.CheckErr(err)
		}

		if changed {
			gitDiff := git.Diff(cfg)
			stacks = stacks.FilterChanged(gitDiff)
		}

		if len(stacks) > 0 {
			fmt.Println("Stacks:")
			for _, stack := range stacks {
				fmt.Printf("- %s\n", stack.Id)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("changed", "c", false, "Filter stacks based on changes made in git")
	listCmd.Flags().BoolP("sorted", "s", false, "Sort listed stacks by order of execution")
}
