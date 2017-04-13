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

func loading() {
	fmt.Println("Creating cluster")
	fmt.Println("This may take a few minutes...")
}

// clusterCmd represents the cluster command
var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "creates a cluster",
	Long:  `creates a cluster in mystack`,
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
		l := log.WithFields(logrus.Fields{
			"controllerURL": config.ControllerURL,
		})
		l.Debug("creating cluster")
		createClusterURL := fmt.Sprintf("%s/clusters/%s/create", config.ControllerURL, clusterName)
		client := models.NewMyStackHTTPClient(config)

		if err != nil {
			l.WithError(err).Fatalf("error during reading file path '%s'", filePath)
		}
		loading()
		body, status, err := client.Put(createClusterURL, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		if status != 200 {
			printer := models.NewErrorPrinter(body, status)
			printer.Print()
			return
		}

		fmt.Printf("Cluster '%s' successfully created\n", clusterName)
	},
}

func init() {
	createCmd.AddCommand(createClusterCmd)
	createClusterCmd.Flags().StringVarP(&clusterName, "clusterName", "c", "", "Name of the cluster to be created")
}
