// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

func deleteCluster(l *logrus.Entry, clusterName string, config *models.Config) {
	l.Debug("deleting cluster")
	deleteClusterURL := fmt.Sprintf("%s/clusters/%s/delete", config.ControllerURL, clusterName)
	client := models.NewMyStackHTTPClient(config)

	fmt.Println("Deleting cluster", clusterName)
	fmt.Println("This may take a few minutes...")

	body, status, err := client.Delete(deleteClusterURL)
	if err != nil {
		l.Fatal(err.Error())
	}

	if status == http.StatusNotFound {
		fmt.Printf("cluster '%s' was not found\n", clusterName)
		return
	} else if status != http.StatusOK {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	fmt.Printf("Cluster '%s' successfully deleted\n", clusterName)
}

// deleteClusterCmd represents the delete_cluster command
var deleteClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deletes a cluster",
	Long:  `deletes a cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		fs := models.NewRealFS()
		c, err := models.ReadConfig(fs, environment)
		if err == nil {
			config = c
		} else {
			log.Fatal("no mystack config file found, you may need to run './mystack login'")
		}

		l := log.WithFields(logrus.Fields{
			"controllerURL": config.ControllerURL,
		})

		if len(args) == 0 {
			fmt.Println("inform cluster name, e.g. './mystack delete cluster mycluster'")
			return
		}

		clusterName := args[0]
		deleteCluster(l, clusterName, config)
	},
}

func init() {
	deleteCmd.AddCommand(deleteClusterCmd)
}
