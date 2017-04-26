// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import "fmt"

//ColumnPrinter implements the Printer interface
type ColumnPrinter struct {
	Title  string
	Column []interface{}
}

//Print formats and prints a JSON
func (j *ColumnPrinter) Print() {
	fmt.Println(j.Title)

	for _, item := range j.Column {
		fmt.Println(item)
	}
}
