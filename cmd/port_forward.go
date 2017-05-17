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
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/mystack-cli/models"
)

var tcpPort int
var interactive bool
var portsFile string

var startPort int = 28000
var wg *sync.WaitGroup

type Port struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

func readPortsFile(portsFile string) (map[string]int, error) {
	portsMap := make(map[string]int)
	if portsFile == "" {
		return portsMap, nil
	}

	ports := struct {
		Ports []*Port `yaml:"ports"`
	}{}

	yamlFile, err := ioutil.ReadFile(portsFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &ports)
	if err != nil {
		return nil, err
	}

	for _, port := range ports.Ports {
		portsMap[port.Name] = port.Port
	}

	return portsMap, nil
}

func get(l *logrus.Entry, clusterName, mystackType string) ([]byte, error) {
	url := fmt.Sprintf("%s/clusters/%s/%s", config.ControllerURL, clusterName, mystackType)
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

	return body, nil
}

func getServiceNames(l *logrus.Entry, clusterName string) ([]string, error) {
	l.Debug("ready to get cluster services")
	body, err := get(l, clusterName, "services")
	if err != nil {
		return nil, err
	}

	bodyJSON := make(map[string][]string)
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	services := bodyJSON["services"]

	return services, nil
}

func getAppNames(l *logrus.Entry, clusterName string) ([]string, error) {
	l.Debug("ready to get cluster apps")
	body, err := get(l, clusterName, "apps")
	if err != nil {
		return nil, err
	}

	bodyJSON := make(map[string]map[string][]string)
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	apps := make([]string, len(bodyJSON["domains"]))
	i := 0
	for name := range bodyJSON["domains"] {
		apps[i] = name
		i = i + 1
	}

	sort.Strings(apps)
	return apps, nil
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

func bindType(
	shiftPort int,
	wg *sync.WaitGroup,
	printer *models.PortForwardPrinter,
	mystackType string,
	l logrus.FieldLogger,
	portsMap map[string]int,
) (int, string) {
	tcpPort := startPort + shiftPort
	if port, contains := portsMap[mystackType]; contains {
		tcpPort = port
	} else if interactive {
		tcpPort, shiftPort = readPort(mystackType, shiftPort)
	} else {
		shiftPort = shiftPort + 1
	}

	localAddr := fmt.Sprintf(":%d", tcpPort)
	remoteAddr := fmt.Sprintf("%s:28000", strings.Split(config.ControllerURL, "://")[1])

	message := map[string]interface{}{
		"token":   config.Token,
		"service": mystackType,
	}

	proxy := models.NewProxy(localAddr, remoteAddr, message)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := proxy.Start()
		if err != nil {
			l.WithError(err).Fatal("error while port forwarding")
		}
	}()

	return shiftPort, localAddr
}

func bindPorts(apps, services []string, l logrus.FieldLogger) error {
	portsMap, err := readPortsFile(portsFile)
	if err != nil {
		return err
	}

	shiftPort := 0
	printer := models.NewPortForwarderPrinter()
	var localAddr string
	for _, service := range services {
		shiftPort, localAddr = bindType(shiftPort, wg, printer, service, l, portsMap)
		printer.AddSvc(service, localAddr)
	}

	for _, app := range apps {
		shiftPort, localAddr = bindType(shiftPort, wg, printer, app, l, portsMap)
		printer.AddApp(app, localAddr)
	}
	printer.Print()

	return nil
}

// portForwardCmd represents the portForward command
var portForwardCmd = &cobra.Command{
	Use:   "port-forward",
	Short: "Forward service ports",
	Long:  `Binds local ports to access mystack services on kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		log := createLog()

		wg = &sync.WaitGroup{}

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

		apps, err := getAppNames(l, args[0])
		if err != nil {
			l.WithError(err).Fatal("error while port forwarding")
		}

		err = bindPorts(apps, services, l)
		if err != nil {
			l.WithError(err).Fatalf("Error reading file '%s'", portsFile)
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(portForwardCmd)
	portForwardCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interative mode to choose ports to bind")
	portForwardCmd.Flags().StringVarP(&portsFile, "ports", "p", "", "ports yaml file with ports biding")
}
