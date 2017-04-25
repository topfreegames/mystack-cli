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
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

var forwardToDNS string
var dnsPort int

type domainList struct {
	Domains []string `json:"domains"`
}

func parseBody(body []byte, status int) (*domainList, error) {
	if status != http.StatusOK {
		errorBody := make(map[string]string)
		err := json.Unmarshal(body, &errorBody)
		if err != nil {
			return nil, fmt.Errorf(string(body))
		}

		return nil, fmt.Errorf("Error: %s\nDescription: %s", errorBody["error"], errorBody["description"])
	}

	bodyMap := make(map[string]map[string][]string)
	err := json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, err
	}

	domainList := &domainList{}

	for _, apps := range bodyMap["domains"] {
		domainList.Domains = append(domainList.Domains, apps...)
	}

	return domainList, nil
}

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "starts a local dns server",
	Long:  `starts a local dns server that will point all stack custom domains to mystack-router, e.g. test production apps`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		if clusterName == "" {
			log.Fatal("inform cluster name with -c flag")
		}

		c, err := models.ReadConfig(environment)
		if err == nil {
			config = c
		} else {
			log.Fatal("no mystack config file found, you may need to run ./mysctl login:", err)
		}
		l := log.WithFields(logrus.Fields{
			"forwardToDNS":  forwardToDNS,
			"controllerURL": config.ControllerURL,
		})
		l.Debug("starting dns server, querying mystack controller for domains and pointTo")
		l.Debug("grabbing mystack custom domain list and router ip from controller")
		getDNSConfigURL := fmt.Sprintf("%s/clusters/%s/apps", config.ControllerURL, clusterName)
		client := models.NewMyStackHTTPClient(config)
		body, status, err := client.Get(getDNSConfigURL)
		if err != nil {
			log.Fatal(err.Error())
		}

		domainList, err := parseBody(body, status)
		if err != nil {
			log.Fatal(err)
		}
		log.WithField("domains-list", domainList).Debug("got domain list from mystack")
		controllerDomain := strings.TrimPrefix(config.ControllerHost, "controller.")
		server, err := models.NewDNSServer(domainList.Domains, forwardToDNS, config.ControllerURL, controllerDomain, dnsPort, log)
		if err != nil {
			log.Fatal(err)
		}
		server.Serve()
	},
}

func init() {
	dnsCmd.Flags().StringVarP(&forwardToDNS, "forwartToDNS", "f", "8.8.8.8:53", "The DNS to forward requests to when not in mystack custom domain list")
	dnsCmd.Flags().IntVarP(&dnsPort, "port", "p", 53, "The port on which UDP Server will listen")
	dnsCmd.Flags().StringVarP(&clusterName, "clusterName", "c", "", "Name of the cluster to be created")
	RootCmd.AddCommand(dnsCmd)
}
