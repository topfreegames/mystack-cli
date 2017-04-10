// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var createConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "creates a cluster config",
	Long:  `creates a cluster config in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	createCmd.AddCommand(createConfigCmd)
}
