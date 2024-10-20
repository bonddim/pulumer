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

	"github.com/bonddim/pulumer/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	pulumer *cli.CLI
	cfgFile string
	dir     string

	rootCmd = &cobra.Command{
		Use:          "pulumer",
		Short:        "Pulumi wrapper to manage multiple local stacks",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error

			if dir != "" {
				if err = os.Chdir(dir); err != nil {
					cobra.CheckErr(err)
				}
			}

			// Bind flags to the viper before creating the CLI
			cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
				if err = viper.BindPFlag(flag.Name, flag); err != nil {
					cobra.CheckErr(err)
				}
			})

			// Create a new CLI
			pulumer, err = cli.NewCLI(cmd.Context(), cfgFile)
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags for all commands
	rootCmd.PersistentFlags().StringVarP(&dir, "chdir", "C", "", "Set working directory")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Set config file (default is .pulumer.yaml in the working directory)")
}

func addChangeDetectionFlags(flags *pflag.FlagSet) {
	flags.BoolP("changed", "c", false, "Filter stacks based on changes made in git")
	flags.String("since", "origin/main", "Set the base ref for the comparison")
}
