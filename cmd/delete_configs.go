// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

// delete_configsCmd represents the delete_configs command
var deleteConfigsCmd = &cobra.Command{
	Use:   "configs CLUSTER_CONFIG [CLUSTER_CONFIGS...]",
	Short: "deletes one or more configs in mystack.",
	Long: `deletes one or more configs in mystack.
CLUSTER_CONFIG is a necessary parameter used to identify specific config.
It's possible to pass one or more config names and they will all be deleted.`,
	Run: DeleteConfigRun,
}

func init() {
	deleteCmd.AddCommand(deleteConfigsCmd)
}
