// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var filePath string
var clusterName string

func getTokenPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/.mystack/token", usr.HomeDir), nil
}

func check(cmdL *logrus.Entry) {
	if filePath == "" {
		cmdL.Fatal("Path to config file must not be empty")
		os.Exit(1)
	}

	if clusterName == "" {
		cmdL.Fatal("Cluster name must not be empty")
		os.Exit(1)
	}
}

func createBody(filePath string) (*io.Reader, error) {
	bts, err := ioutil.ReadAll(filePath)
	if err != nil {
		return nil, err
	}

	return nil
}

// configCmd represents the config command
var createConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "creates a cluster config",
	Long:  `creates a cluster config in mystack`,
	Run: func(cmd *cobra.Command, args []string) {
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
			"source":    "createConfigCmd",
			"operation": "Run",
		})

		check(cmdL)

		client := &http.Client{}
		file
		if err != nil {
			msg := fmt.Sprintf("Failed to open config file '%s'", filePath)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		createClusterURL := fmt.Sprintf("%s/cluster-configs/%s/create", controllerURL, clusterName)
		request, err := http.NewRequest("PUT", createClusterURL, file)
		if err != nil {
			msg := fmt.Sprintf("Failed to build request to '%s'", controllerURL)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		tokenPath, err := getTokenPath()
		if err != nil {
			msg := "Failed to get Access Token (Have you logged in?)"
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		token, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			msg := "Failed to read Access Token file"
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		auth := fmt.Sprintf("Bearer %s", string(token))
		request.Header.Add("Authorization", auth)

		response, err := client.Do(request)
		if err != nil {
			msg := fmt.Sprintf("Failed to execute request to '%s'", controllerURL)
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			msg := "Failed to read response"
			cmdL.WithError(err).Fatal(msg)
			os.Exit(1)
		}

		bodyJSON := make(map[string]interface{})
		json.Unmarshal(body, &bodyJSON)
		if response.StatusCode != 200 && response.StatusCode != 201 {
			fmt.Printf(
				"Status: %d\nError: %s\nDescription: %s",
				response.StatusCode,
				bodyJSON["error"],
				bodyJSON["description"],
			)
		}
	},
}

func init() {
	createCmd.AddCommand(createConfigCmd)
	createConfigCmd.Flags().StringVarP(&controllerURL, "controllerURL", "s", "", "Controllers URL")
	createConfigCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "Config file path")
	createConfigCmd.Flags().StringVarP(&clusterName, "clusterName", "c", "", "Name of the cluster to be created")
}
