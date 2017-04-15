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

// getConfigCmd represents the get_config command
var getConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "list or get cluster configs",
	Long:  `list or get a cluster configs from mystack`,
	Run: func(cmd *cobra.Command, args []string) {
		ll := logrus.InfoLevel
		switch verbose {
		case 0:
			ll = logrus.ErrorLevel
			break
		case 1:
			ll = logrus.WarnLevel
			break
		case 3:
			ll = logrus.DebugLevel
			break
		default:
			ll = logrus.InfoLevel
		}

		log = logrus.New()
		log.Level = ll

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
			"controllerURL": config.ControllerURL,
		})
		l.Debug("ready to get cluster config")
		url := fmt.Sprintf("%s/cluster-configs/%s", config.ControllerURL, clusterName)
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

		bodyJSON := make(map[string]interface{})
		json.Unmarshal(body, &bodyJSON)
		printer := &models.ObjPrinter{
			Title: "CLUSTER-CONFIG",
			Obj:   bodyJSON["yaml"],
		}
		printer.Print()
	},
}

func init() {
	getCmd.AddCommand(getConfigCmd)
	getConfigCmd.Flags().StringVarP(&clusterName, "clusterName", "c", "", "Cluster name to get config")
	getClusterCmd.MarkFlagRequired("clusterName")
}
