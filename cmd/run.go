// /*
// Copyright © 2024 Dmytro Bondar <git@bonddim.com>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
// */
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	login bool
	yes   bool

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run pulumi command",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			// Check for destroy command confirmation
			if args[0] == "destroy" && !yes {
				cobra.CheckErr("use --yes flag to confirm destroy")
			}

			if login {
				if os.Getenv("PULUMI_BACKEND_URL") == "" && os.Getenv("PULUMI_ACCESS_TOKEN") == "" {
					cobra.CheckErr("PULUMI_BACKEND_URL or PULUMI_ACCESS_TOKEN environment variable is required for login.")
				} else if err := pulumer.Run([]string{"login"}); err != nil {
					cobra.CheckErr(err)
				}
			}
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			return pulumer.Run(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	addChangeDetectionFlags(runCmd.Flags())
	runCmd.Flags().BoolVar(&login, "login", false, "Login to pulumi backend. Requires PULUMI_BACKEND_URL or PULUMI_ACCESS_TOKEN environment variable")
	runCmd.Flags().BoolVar(&yes, "yes", false, "Extra confirmation for destroy command")
}
