// mystack
// https://github.com/topfreegames/mystack/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Verbose determines how verbose mystack wll run under
var Verbose int

// RootCmd is the root command for mystack CLI application
var RootCmd = &cobra.Command{
	Use:   "mystack",
	Short: "mystack handles manages your personal cluster",
	Long:  `Use mystack to start your services on kubernetes.`,
}

// Execute runs RootCmd to initialize mystack CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(
		&Verbose, "verbose", "v", 0,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug",
	)
}
