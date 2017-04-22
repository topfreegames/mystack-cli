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

// listConfigsCmd represents the listConfigs.go command
var listConfigsCmd = &cobra.Command{
	Use:   "list",
	Short: "list cluster configs",
	Long:  `Get the list of cluster configs already created on the Mystack-Controller`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			log.Fatal("no mystack config file found, you may need to run ./mysctl login")
		}
		l := log.WithFields(logrus.Fields{
			"controllerURL": config.ControllerURL,
		})
		l.Debug("ready to get cluster config list")
		url := fmt.Sprintf("%s/cluster-configs", config.ControllerURL)
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
		printer := &models.ColumnPrinter{
			Title:  "CLUSTER-NAMES",
			Column: bodyJSON["names"],
		}
		printer.Print()
	},
}

func init() {
	getCmd.AddCommand(listConfigsCmd)
}
