// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// updateCm represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates a resource",
	Long:  `updates a resource in mystack`,
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
