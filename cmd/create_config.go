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

func createBody() (map[string]interface{}, error) {
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(bts) == 0 {
		return nil, fmt.Errorf("file path was not informed")
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
		log := createLog()

		cmdL := log.WithFields(logrus.Fields{
			"source":    "createConfigCmd",
			"operation": "Run",
		})

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			cmdL.WithError(err).Fatal("no mystack config file found, you may need to run ./mystack login")
		}

		client := models.NewMyStackHTTPClient(config)
		createClusterURL := fmt.Sprintf("%s/cluster-configs/%s/create", c.ControllerURL, clusterName)
		bodyJSON, err := createBody()
		if err != nil {
			cmdL.WithError(err).Fatalf("error during reading file path '%s'", filePath)
		}

		body, status, err := client.Put(createClusterURL, bodyJSON)
		if err != nil {
			msg := fmt.Sprintf("Failed to execute request to '%s'", c.ControllerURL)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		if status != 200 && status != 201 {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		fmt.Printf("Cluster config '%s' successfully created\n", clusterName)
	},
}

func init() {
	createCmd.AddCommand(createConfigCmd)
	createConfigCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Config file path")
}
