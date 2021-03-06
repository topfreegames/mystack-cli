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

var loginCmd = &cobra.Command{
	Use:   "login CONTROLLER_URL",
	Short: "login on mystack",
	Long: `First login on mystack to get access on your personal stack of services running on Kubernetes.
CONTROLLER_URL is a necessary parameter used to connect to Mystack-Controller.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		cmdL := log.WithFields(logrus.Fields{
			"source":    "loginCmd",
			"operation": "Run",
			"debug":     debug,
		})

		if len(args) == 0 {
			cmdL.Fatal("inform controller url")
		}

		controllerURL = args[0]

		cmdL.Debug("Creating callback server...")
		app, err := api.NewApp(
			host,
			port,
			debug,
			log,
			environment,
			controllerURL,
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
	loginCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	loginCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (log level error)")
}
