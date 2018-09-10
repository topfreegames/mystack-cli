// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

// updateConfigCmd represents the config command
var updateConfigCmd = &cobra.Command{
	Use:   "config CLUSTER_CONFIG",
	Short: "updates a cluster config",
	Long: `updates a cluster config in mystack.
CLUSTER_CONFIG is a necessary parameter used to name cluster config.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		cmdL := log.WithFields(logrus.Fields{
			"source":    "updateConfigCmd",
			"operation": "Run",
		})

		fs := models.NewRealFS()
		c, err := models.ReadConfig(fs, environment)
		if err == nil {
			config = c
		} else {
			fmt.Println("no mystack config file found, you may need to run './mystack login'")
			return
		}

		if filePath == "" {
			fmt.Println("inform config file, e.g. './mystack update config myconfig -f /path/to/config/file'")
			return
		}
		if len(args) == 0 {
			fmt.Println("inform cluster name, e.g. './mystack update config myconfig'")
			return
		}

		clusterName := args[0]

		client := models.NewMyStackHTTPClient(config)
		updateClusterURL := fmt.Sprintf(
			"%s/cluster-configs/%s/update",
			config.ControllerURL,
			clusterName,
		)
		bodyJSON, err := createBody()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}

		body, status, err := client.Put(updateClusterURL, bodyJSON)
		if err != nil {
			msg := fmt.Sprintf("Failed to execute request to '%s'", config.ControllerURL)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		if status == http.StatusConflict {
			title := fmt.Sprintf("config '%s' already exists", clusterName)
			msg := fmt.Sprintf("to update it, you have to run './mystack delete config %s' and update a new one", clusterName)
			printer := models.NewStrLogPrinter(msg, title)
			printer.Print()
			return
		} else if status != http.StatusOK {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		fmt.Printf("Cluster config '%s' successfully updated\n", clusterName)
	},
}

func init() {
	updateCmd.AddCommand(updateConfigCmd)
	updateConfigCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Config file path")
}
