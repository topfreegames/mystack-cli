// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/api"
)

const host string = "0.0.0.0"
const port int = 57459

var debug bool
var quiet bool
var controllerURL string
var controllerHost string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login on mystack",
	Long:  "First login on mystack to get access on your personal stack of services running on Kubernetes",
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

		cmdL := log.WithFields(logrus.Fields{
			"source":    "loginCmd",
			"operation": "Run",
			"debug":     debug,
		})

		if len(controllerHost) == 0 {
			controllerHost = controllerURL
		}

		cmdL.Debug("Creating callback server...")
		app, err := api.NewApp(
			host,
			port,
			debug,
			log,
			environment,
			controllerURL,
			controllerHost,
		)
		if err != nil {
			cmdL.WithError(err).Fatal("Failed to start server.")
		}
		cmdL.Debug("Application created successfully.")

		cmdL.Debug("Starting application...")
		closer, err := app.ListenAndLoginAndServe()
		if closer != nil {
			defer closer.Close()
		}
		if err != nil {
			cmdL.WithError(err).Fatal("Error running application.")
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&controllerURL, "controllerURL", "s", "http://localhost:80", "Controllers URL")
	loginCmd.Flags().StringVarP(&controllerHost, "controllerHost", "o", "", "Controller Host")
	loginCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	loginCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (log level error)")
}
