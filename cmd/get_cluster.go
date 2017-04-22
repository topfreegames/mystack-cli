// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

// getClusterCmd represents the get_cluster command
var getClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "list or get clusters",
	Long:  `list or get cluster in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			log.Fatal("no mystack config file found, you may need to run ./mysctl login")
		}

		if len(clusterName) == 0 {
			log.Fatal("clusterName must be informed with flag -c")
		}
		l := log.WithFields(logrus.Fields{
			"controllerURL":  config.ControllerURL,
			"controllerHost": config.ControllerHost,
		})
		l.Debug("ready to get cluster config")
		url := fmt.Sprintf("%s/clusters/%s/apps", config.ControllerURL, clusterName)
		client := models.NewMyStackHTTPClient(config)
		body, status, err := client.Get(url)
		if err != nil {
			log.Fatal(err.Error())
		}

		if status != 200 {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		bodyJSON := make(map[string][]interface{})
		json.Unmarshal(body, &bodyJSON)
		printer := &models.RoutePrinter{
			Domain: config.ControllerHost,
			Apps:   bodyJSON["apps"],
		}
		printer.Print()
	},
}

func init() {
	getCmd.AddCommand(getClusterCmd)
}
