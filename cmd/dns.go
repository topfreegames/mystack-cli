// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

var forwardToDNS string
var dnsPort int

type domainList struct {
	Domains []string `json:"domains"`
}

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "starts a local dns server",
	Long:  `starts a local dns server that will point all stack custom domains to mystack-router, e.g. test production apps`,
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
			"forwardToDNS":  forwardToDNS,
			"controllerURL": config.ControllerURL,
		})
		l.Debug("starting dns server, querying mystack controller for domains and pointTo")
		l.Debug("grabbing mystack custom domain list and router ip from controller")
		getDNSConfigURL := fmt.Sprintf("%s/dns", config.ControllerURL)
		client := models.NewMyStackHTTPClient(config)
		body, err := client.Get(getDNSConfigURL)
		if err != nil {
			log.Fatal(err.Error())
		}
		domainList := &domainList{}
		json.Unmarshal(body, domainList)
		log.WithField("domains-list", domainList).Debug("got domain list from mystack")
		server, err := models.NewDNSServer(domainList.Domains, forwardToDNS, config.ControllerURL, dnsPort, log)
		if err != nil {
			log.Fatal(err)
		}
		server.Serve()
	},
}

func init() {
	dnsCmd.Flags().StringVarP(&forwardToDNS, "forwartToDNS", "f", "8.8.8.8:53", "The DNS to forward requests to when not in mystack custom domain list")
	dnsCmd.Flags().IntVarP(&dnsPort, "port", "p", 53, "The port on which UDP Server will listen")
	RootCmd.AddCommand(dnsCmd)
}
