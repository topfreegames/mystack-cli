// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/metadata"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays mystack cli version",
	Long:  `displays mystack cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", metadata.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
