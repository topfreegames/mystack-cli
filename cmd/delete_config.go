// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

// delete_configCmd represents the delete_config command
var deleteConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "deletes a config",
	Long:  `deletes a config in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		ll := logrus.InfoLevel
		switch verbose {
		case 0:
			ll = logrus.ErrorLevel
		case 1:
			ll = logrus.WarnLevel
		case 3:
			ll = logrus.DebugLevel
		}

		var log = logrus.New()
		log.Formatter = new(logrus.JSONFormatter)
		log.Level = ll

		cmdL := log.WithFields(logrus.Fields{
			"source":    "createConfigCmd",
			"operation": "Run",
		})

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			cmdL.WithError(err).Fatal("no mystack config file found, you may need to run ./mysctl login")
		}

		client := models.NewMyStackHTTPClient(config)
		deleteClusterURL := fmt.Sprintf("%s/cluster-configs/%s/remove", c.ControllerURL, clusterName)
		if err != nil {
			cmdL.WithError(err).Fatalf("error during reading file path '%s'", filePath)
		}

		body, status, err := client.Delete(deleteClusterURL)
		if err != nil {
			msg := fmt.Sprintf("Failed to execute request to '%s'", c.ControllerURL)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		if status != 200 && status != 201 {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		fmt.Printf("Cluster config '%s' successfully deleted\n", clusterName)
	},
}

func init() {
	deleteCmd.AddCommand(deleteConfigCmd)
	deleteConfigCmd.Flags().StringVarP(&clusterName, "clusterName", "c", "", "Name of the cluster to be created")
}
