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

// getConfigsCmd represents the configs command
var getConfigsCmd = &cobra.Command{
	Use:   "configs [CLUSTER_CONFIG...]",
	Short: "list or get cluster configs",
	Long: `list or get cluster configs from mystack.
CLUSTER_CONFIG is used to fetch specific cluster.
It's possible to pass one or more cluster names and they will be all fetched.
If no CLUSTER_CONFIG is passed, a list with the saved cluster configs is returned.`,
	Run: GetConfigRun,
}

func init() {
	getCmd.AddCommand(getConfigsCmd)
}
