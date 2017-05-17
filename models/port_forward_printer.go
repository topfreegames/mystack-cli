// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

//PortForwardPrinter implements the Printer interface
//It prints the port binded for each mystack service of a cluster
type PortForwardPrinter struct {
	svcs map[string]string
	apps map[string]string
}

//NewPortForwarderPrinter returns a PortForwardPrinter
func NewPortForwarderPrinter() *PortForwardPrinter {
	return &PortForwardPrinter{
		svcs: make(map[string]string),
		apps: make(map[string]string),
	}
}

//Add a new port binded to a service
func (p *PortForwardPrinter) AddSvc(service, host string) {
	p.svcs[service] = host
}

//Add a new port binded to a app
func (p *PortForwardPrinter) AddApp(app, host string) {
	p.apps[app] = host
}

//Print formats and prints the ports for each service
func (p *PortForwardPrinter) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	services := make([]string, len(p.svcs))
	i := 0
	for name, _ := range p.svcs {
		services[i] = name
		i = i + 1
	}
	sort.Strings(services)

	apps := make([]string, len(p.apps))
	i = 0
	for name, _ := range p.apps {
		apps[i] = name
		i = i + 1
	}
	sort.Strings(apps)

	fmt.Fprintf(w, "%s \t%v\n", "SERVICES", "LOCAL HOSTNAME")
	for _, name := range services {
		fmt.Fprintf(w, "%s \t%v\n", name, p.svcs[name])
	}
	fmt.Fprintf(w, "\n%s \t%v\n", "APPS", "LOCAL HOSTNAME")
	for _, name := range apps {
		fmt.Fprintf(w, "%s \t%v\n", name, p.apps[name])
	}

	w.Flush()
}
