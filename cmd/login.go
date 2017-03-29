package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack/mystack-cli/api"
)

const host string = "0.0.0.0"
const port int = 8989

var debug bool
var quiet bool

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login on mystack",
	Long:  "First login on mystack to get access on your personal stack of services running on Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig()
		ll := logrus.InfoLevel
		switch Verbose {
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
			"source":    "loginCmd",
			"operation": "Run",
			"debug":     debug,
		})

		cmdL.Debug("Creating callback server...")
		app, err := api.NewApp(
			host,
			port,
			debug,
			log,
			config,
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
