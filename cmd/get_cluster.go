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

func getCluster(l *logrus.Entry, clusterName string, config *models.Config) {
	l.Debug("ready to get cluster config")
	url := fmt.Sprintf("%s/clusters/%s/apps", config.ControllerURL, clusterName)
	client := models.NewMyStackHTTPClient(config)
	body, status, err := client.Get(url, config.ControllerHost)
	if err != nil {
		log.Fatal(err.Error())
	}

	if status == http.StatusNotFound {
		title := fmt.Sprintf("cluster '%s' was not found", clusterName)
		msg := fmt.Sprintf("you have to run './mystack create cluster %s'", clusterName)
		printer := models.NewStrLogPrinter(msg, title)
		printer.Print()
		return
	} else if status != http.StatusOK {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	bodyJSON := make(map[string]map[string][]string)
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		log.Fatal(err.Error())
	}

	printer := &models.RoutePrinter{
		Domain: config.ControllerHost,
		Apps:   bodyJSON["domains"],
	}
	printer.Print()
}

// getClusterCmd represents the get_cluster command
var getClusterCmd = &cobra.Command{
	Use:   "cluster CLUSTER_NAME",
	Short: "list or get clusters",
	Long: `list or get cluster in mystack.
CLUSTER_NAME is a necessary parameter used to fetch specific cluster.`,
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
			"controllerURL":  config.ControllerURL,
			"controllerHost": config.ControllerHost,
		})

		if len(args) == 0 {
			fmt.Println("cluster name must be informed, e.g './mystack get cluster mycluster'")
			return
		}

		getCluster(l, args[0], c)
	},
}

func init() {
	getCmd.AddCommand(getClusterCmd)
}
