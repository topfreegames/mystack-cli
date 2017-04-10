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

// clusterCmd represents the cluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "creates a cluster",
	Long:  `creates a cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cluster called")
	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)
}
