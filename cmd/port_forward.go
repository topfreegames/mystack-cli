// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

var tcpPort int
var interactive bool

var startPort int = 28000
var wg sync.WaitGroup

func getServiceNames(l *logrus.Entry, clusterName string) ([]string, error) {
	l.Debug("ready to get cluster services")
	url := fmt.Sprintf("%s/clusters/%s/services", config.ControllerURL, clusterName)
	client := models.NewMyStackHTTPClient(config)
	body, status, err := client.Get(url, config.ControllerHost)
	if err != nil {
		return nil, err
	}

	if status == http.StatusNotFound {
		return nil, fmt.Errorf("cluster '%s' was not found\nyou have to run './mystack create cluster %s'", clusterName, clusterName)
	} else if status != http.StatusOK {
		return nil, fmt.Errorf("Status: %d\nBody: %s", status, body)
	}

	bodyJSON := make(map[string][]string)
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return bodyJSON["services"], nil
}

func readPort(service string, shiftPort int) (int, int) {
	tcpPort = startPort + shiftPort
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Choose port to bind for %s [Default %d]: ", service, tcpPort)
		portStr, _ := reader.ReadString('\n')
		portStr = strings.TrimSpace(portStr)

		if portStr == "" {
			shiftPort = shiftPort + 1
			break
		}

		port64, err := strconv.ParseUint(portStr, 10, 64)
		if err != nil {
			fmt.Println("Error: choose a valid port number")
			continue
		}

		tcpPort = int(port64)
		break
	}

	return tcpPort, shiftPort
}

func bindPorts(services []string, l logrus.FieldLogger) {
	shiftPort := 0
	printer := models.NewPortForwarderPrinter()
	for _, service := range services {
		tcpPort := startPort + shiftPort
		if interactive {
			tcpPort, shiftPort = readPort(service, shiftPort)
		} else {
			shiftPort = shiftPort + 1
		}

		localAddr := fmt.Sprintf(":%d", tcpPort)
		remoteAddr := fmt.Sprintf("%s:28000", strings.Split(config.ControllerURL, "://")[1])

		message := map[string]interface{}{
			"token":   config.Token,
			"service": service,
		}

		proxy := models.NewProxy(localAddr, remoteAddr, message)

		printer.Add(service, localAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := proxy.Start()
			if err != nil {
				l.WithError(err).Fatal("error while port forwarding")
			}
		}()
	}

	printer.Print()
}

// portForwardCmd represents the portForward command
var portForwardCmd = &cobra.Command{
	Use:   "port-forward",
	Short: "Forward service ports",
	Long:  `Binds local ports to access mystack services on kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		fs := models.NewRealFS()
		c, err := models.ReadConfig(fs, environment)
		if err == nil {
			config = c
		} else {
			log.WithError(err).Fatal("no mystack config file found, you may need to run './mystack login'")
		}

		l := log.WithFields(logrus.Fields{
			"controllerURL": config.ControllerURL,
		})

		if len(args) == 0 {
			fmt.Println("cluster name must be informed, e.g './mystack port-forward mycluster'")
			return
		}

		services, err := getServiceNames(l, args[0])
		if err != nil {
			l.WithError(err).Fatal("error while port forwarding")
		}

		bindPorts(services, l)
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(portForwardCmd)
	portForwardCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interative mode to choose ports to bind")
}
