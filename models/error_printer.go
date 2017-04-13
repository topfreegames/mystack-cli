// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

//ErrorPrinter implements the Printer interface
type ErrorPrinter struct {
	Error  map[string]interface{}
	Status int
}

//NewErrorPrinter is the ErrorPrinter ctor
func NewErrorPrinter(body []byte, status int) *ErrorPrinter {
	bodyJSON := make(map[string]interface{})
	json.Unmarshal(body, &bodyJSON)
	return &ErrorPrinter{
		Error:  bodyJSON,
		Status: status,
	}
}

//Print formats and prints a JSON
func (e *ErrorPrinter) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	description := e.Error["description"].(string)

	fmt.Fprintln(w, "An error ocurred:")
	fmt.Fprintf(w, "\tStatus:\t%d\n", e.Status)
	fmt.Fprintf(w, "\tError:\t%s\n", e.Error["error"])
	fmt.Fprintf(w, "\tDescription:\t%s\n", strings.Replace(description, "\n", " ", -1))
	w.Flush()
}
