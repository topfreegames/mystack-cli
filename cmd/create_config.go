// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

var filePath string
var clusterName string

func createBody() (map[string]interface{}, error) {
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	bodyJSON := make(map[string]interface{})
	bodyJSON["yaml"] = string(bts)

	return bodyJSON, nil
}

// configCmd represents the config command
var createConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "creates a cluster config",
	Long:  `creates a cluster config in mystack`,
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
		createClusterURL := fmt.Sprintf("%s/cluster-configs/%s/create", controllerURL, clusterName)
		bodyJSON, err := createBody()
		if err != nil {
			cmdL.WithError(err).Fatalf("error during reading file path '%s'", filePath)
		}

		body, status, err := client.Put(createClusterURL, bodyJSON)
		if err != nil {
			msg := fmt.Sprintf("Failed to execute request to '%s'", controllerURL)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		if status != 200 && status != 201 {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		fmt.Println("Success")
	},
}

func init() {
	createCmd.AddCommand(createConfigCmd)
	createConfigCmd.Flags().StringVarP(&controllerURL, "controllerURL", "s", "", "Controllers URL")
	createConfigCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Config file path")
	createConfigCmd.Flags().StringVarP(&clusterName, "clusterName", "c", "", "Name of the cluster to be created")
}
