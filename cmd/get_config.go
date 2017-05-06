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

func getConfig(l *logrus.Entry, clusterName string, config *models.Config) {
	l.Debugf("ready to get %s cluster config", clusterName)
	url := fmt.Sprintf("%s/cluster-configs/%s", config.ControllerURL, clusterName)
	client := models.NewMyStackHTTPClient(config)
	body, status, err := client.Get(url, config.ControllerHost)
	if err != nil {
		l.Fatal(err.Error())
	}

	if status == http.StatusNotFound {
		title := fmt.Sprintf("config '%s' was not found", clusterName)
		msg := fmt.Sprintf("you have to run './mystack create config %s'", clusterName)
		printer := models.NewStrLogPrinter(msg, title)
		printer.Print()
		return
	} else if status != http.StatusOK {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	bodyJSON := make(map[string]interface{})
	json.Unmarshal(body, &bodyJSON)
	printer := models.NewStrLogPrinter(bodyJSON["yaml"].(string), fmt.Sprintf("config '%s'", clusterName))
	printer.Print()
}

func listConfigs(l *logrus.Entry, config *models.Config) {
	url := fmt.Sprintf("%s/cluster-configs", config.ControllerURL)
	client := models.NewMyStackHTTPClient(config)
	body, status, err := client.Get(url, config.ControllerHost)
	if err != nil {
		log.Fatal(err.Error())
	}

	if status != http.StatusOK {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	bodyJSON := make(map[string][]interface{})
	json.Unmarshal(body, &bodyJSON)

	if len(bodyJSON["names"]) == 0 {
		fmt.Println("There is no cluster configs. Run 'mystack create config' to add a new one")
		return
	}

	printer := &models.ColumnPrinter{
		Title:  "CLUSTER-NAMES",
		Column: bodyJSON["names"],
	}
	printer.Print()
}

func GetConfigRun(cmd *cobra.Command, args []string) {
	log := createLog()

	c, err := models.ReadConfig(environment)
	if err == nil {
		config = c
	} else {
		log.Fatal("no mystack config file found, you may need to run './mystack login'")
	}

	l := log.WithFields(logrus.Fields{
		"controllerURL":  config.ControllerURL,
		"controllerHost": config.ControllerHost,
		"args":           args,
	})

	if len(args) == 0 {
		listConfigs(l, config)
		return
	}

	for i, clusterName := range args {
		getConfig(l, clusterName, config)
		if i < len(args)-1 {
			fmt.Println("")
		}
	}
}

// getConfigCmd represents the get_config command
var getConfigCmd = &cobra.Command{
	Use:   "config [CLUSTER_CONFIG...]",
	Short: "list or get cluster configs",
	Long: `list or get cluster configs from mystack.
CLUSTER_CONFIG is used to fetch specific cluster.
It's possible to pass one or more cluster configs and they will be all fetched.
If no CLUSTER_CONFIG is passed, a list with the saved cluster configs is returned.`,
	Run: GetConfigRun,
}

func init() {
	getCmd.AddCommand(getConfigCmd)
}
