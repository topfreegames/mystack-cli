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

// getConfigCmd represents the get_config command
var getConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "list or get cluster configs",
	Long:  `list or get a cluster configs from mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get_config called")
	},
}

func init() {
	getCmd.AddCommand(getConfigCmd)
}
