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

// deleteClusterCmd represents the delete_cluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deletes a cluster",
	Long:  `deletes a cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete_cluster called")
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)
}
