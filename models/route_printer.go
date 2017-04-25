// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
	"strings"
)

//RoutePrinter implements the Printer interface
type RoutePrinter struct {
	Apps   []interface{}
	Domain string
}

//Print formats and prints a JSON
func (j *RoutePrinter) Print() {
	fmt.Println("APP ROUTES")

	domain := strings.TrimPrefix(j.Domain, "controller.")
	for _, item := range j.Apps {
		fmt.Printf("%s.%s\n", item, domain)
	}
}
