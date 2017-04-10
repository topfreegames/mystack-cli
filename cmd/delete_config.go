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

// delete_configCmd represents the delete_config command
var deleteConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "deletes a config",
	Long:  `deletes a config in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete_config called")
	},
}

func init() {
	deleteCmd.AddCommand(deleteConfigCmd)
}
