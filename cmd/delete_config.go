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

func deleteConfig(l *logrus.Entry, clusterName string, config *models.Config) {
	client := models.NewMyStackHTTPClient(config)
	deleteClusterURL := fmt.Sprintf(
		"%s/cluster-configs/%s/remove",
		config.ControllerURL,
		clusterName,
	)

	body, status, err := client.Delete(deleteClusterURL)
	if err != nil {
		msg := fmt.Sprintf("Failed to execute request to '%s'", config.ControllerURL)
		l.WithError(err).Fatal(msg)
		os.Exit(1)
	}

	if status != 200 && status != 201 {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	fmt.Printf("Cluster config '%s' successfully deleted\n", clusterName)
}

func DeleteConfigRun(cmd *cobra.Command, args []string) {
	log := createLog()

	l := log.WithFields(logrus.Fields{
		"source":    "createConfigCmd",
		"operation": "Run",
	})

	c, err := models.ReadConfig(environment)
	if err == nil {
		config = c
	} else {
		l.WithError(err).Fatal("no mystack config file found, you may need to run './mystack login'")
	}

	if len(args) == 0 {
		fmt.Println("inform cluster name, e.g. './mystack delete config myconfig'")
		return
	}

	for _, clusterName := range args {
		deleteConfig(l, clusterName, config)
	}
}

// delete_configCmd represents the delete_config command
var deleteConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "deletes a config",
	Long:  `deletes a config in mystack`,
	Run:   DeleteConfigRun,
}

func init() {
	deleteCmd.AddCommand(deleteConfigCmd)
}
