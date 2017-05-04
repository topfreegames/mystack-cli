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

//RoutePrinter implements the Printer interface
type RoutePrinter struct {
	Apps   map[string][]string
	Domain string
}

//Print formats and prints a JSON
func (j *RoutePrinter) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	names := make([]string, len(j.Apps))
	i := 0
	for name, _ := range j.Apps {
		names[i] = name
		i = i + 1
	}

	sort.Strings(names)

	fmt.Fprintln(w, "APP ROUTES")
	for _, name := range names {
		domains := j.Apps[name]
		fmt.Fprintf(w, "%s:\t%v\n", name, domains)
	}

	w.Flush()
}
