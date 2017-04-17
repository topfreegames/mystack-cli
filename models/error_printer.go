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
	ErrorMap map[string]interface{}
	Status   int
	Error    string
}

//NewErrorPrinter is the ErrorPrinter ctor
func NewErrorPrinter(body []byte, status int) *ErrorPrinter {
	bodyJSON := make(map[string]interface{})
	json.Unmarshal(body, &bodyJSON)
	return &ErrorPrinter{
		ErrorMap: bodyJSON,
		Status:   status,
		Error:    string(body),
	}
}

//Print formats and prints a JSON
func (e *ErrorPrinter) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	if _, ok := e.ErrorMap["description"]; ok {
		description := e.ErrorMap["description"].(string)

		fmt.Fprintln(w, "An error ocurred:")
		fmt.Fprintf(w, "\tStatus:\t%d\n", e.Status)
		fmt.Fprintf(w, "\tError:\t%s\n", e.ErrorMap["error"])
		fmt.Fprintf(w, "\tDescription:\t%s\n", strings.Replace(description, "\n", " ", -1))
		w.Flush()
		return
	}

	fmt.Fprintln(w, e.Error)
}
