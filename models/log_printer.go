// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
)

//LogPrinter implements the Printer interface
type LogPrinter struct {
	Message string
	Title   string
}

//NewLogPrinter is the LogPrinter ctor
func NewLogPrinter(body []byte, title string) *LogPrinter {
	return &LogPrinter{
		Message: string(body),
		Title:   title,
	}
}

//NewLogPrinter is the LogPrinter ctor
func NewStrLogPrinter(body string, title string) *LogPrinter {
	return &LogPrinter{
		Message: body,
		Title:   title,
	}
}

//Print prints log
func (l *LogPrinter) Print() {
	fmt.Println(l.Title)
	fmt.Println(l.Message)
}
