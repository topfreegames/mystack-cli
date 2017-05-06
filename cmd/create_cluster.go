// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

func loading() {
	fmt.Println("Creating cluster")
	fmt.Println("This may take a few minutes...")
}

func createCluster(l *logrus.Entry, clusterName string, config *models.Config) {
	l.Debug("creating cluster")
	createClusterURL := fmt.Sprintf("%s/clusters/%s/create", config.ControllerURL, clusterName)
	client := models.NewMyStackHTTPClient(config)

	loading()
	body, status, err := client.Put(createClusterURL, nil)
	if err != nil {
		l.Fatal(err.Error())
	}

	if status == http.StatusNotFound {
		title := fmt.Sprintf("config '%s' was not found", clusterName)
		msg := fmt.Sprintf("you have to run './mystack create config %s'", clusterName)
		printer := models.NewStrLogPrinter(msg, title)
		printer.Print()
		return
	} else if status == http.StatusConflict {
		fmt.Printf("cluster '%s' already exists\n", clusterName)
		return
	}
	if status != http.StatusOK {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	fmt.Printf("Cluster '%s' successfully created\n", clusterName)

	bodyJSON := make(map[string]map[string][]string)
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		l.Fatal(err.Error())
	}

	printer := &models.RoutePrinter{
		Domain: config.ControllerHost,
		Apps:   bodyJSON["domains"],
	}
	printer.Print()
}

// clusterCmd represents the cluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "creates a cluster",
	Long:  `creates a cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			log.WithError(err).Fatal("no mystack config file found, you may need to run './mystack login'")
		}

		l := log.WithFields(logrus.Fields{
			"controllerURL": config.ControllerURL,
		})

		if len(args) == 0 {
			fmt.Println("inform cluster name, e.g. './mystack create cluster mycluster'")
			return
		}

		clusterName := args[0]
		createCluster(l, clusterName, config)
	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)
}
