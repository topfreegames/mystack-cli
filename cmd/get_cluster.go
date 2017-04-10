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

// getClusterCmd represents the get_cluster command
var getClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "list or get clusters",
	Long:  `list or get cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("get_cluster called")
	},
}

func init() {
	getCmd.AddCommand(getClusterCmd)
}
