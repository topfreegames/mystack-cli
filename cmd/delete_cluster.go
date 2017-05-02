// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

// deleteClusterCmd represents the delete_cluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deletes a cluster",
	Long:  `deletes a cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			log.Fatal("no mystack config file found, you may need to run ./mystack login")
		}
		l := log.WithFields(logrus.Fields{
			"controllerURL": config.ControllerURL,
		})
		l.Debug("deleting cluster")
		createClusterURL := fmt.Sprintf("%s/clusters/%s/delete", config.ControllerURL, clusterName)
		client := models.NewMyStackHTTPClient(config)

		if err != nil {
			l.WithError(err).Fatalf("error during reading file path '%s'", filePath)
		}
		fmt.Println("Deleting cluster")
		fmt.Println("This may take a few minutes...")
		body, status, err := client.Delete(createClusterURL)
		if err != nil {
			log.Fatal(err.Error())
		}

		if status != 200 {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		fmt.Printf("Cluster '%s' successfully deleted\n", clusterName)
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)
}
