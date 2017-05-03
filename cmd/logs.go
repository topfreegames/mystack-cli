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

func getLog(l *logrus.Entry, app string, config *models.Config) {
	l.Debug("ready to get app logs")
	url := fmt.Sprintf("%s/logs/apps/%s", config.ControllerURL, app)
	client := models.NewMyStackHTTPClient(config)
	body, status, err := client.Get(url, config.LoggerHost)
	if err != nil {
		log.Fatal(err.Error())
	}

	if status != 200 {
		printer := models.NewErrorPrinter(body, status)
		printer.Print()
		return
	}

	printer := models.NewLogPrinter(body, app)
	printer.Print()
}

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get app logs",
	Long:  `Get apps' logs running o Mystack user's cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			"loggerHost":     config.LoggerHost,
		})

		if len(args) == 0 {
			fmt.Println("inform app name, e.g. './mystack logs myapp'")
			return
		}

		for i, app := range args {
			getLog(l, app, config)
			if i < len(args)-1 {
				fmt.Println("")
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(logsCmd)
}
