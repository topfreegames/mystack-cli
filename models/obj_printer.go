// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import "fmt"

//ObjPrinter implements the Printer interface
type ObjPrinter struct {
	Title       string
	ClusterName string
	Obj         interface{}
}

//Print formats and prints a JSON
func (j *ObjPrinter) Print() {
	fmt.Println(j.Title)
	fmt.Printf("Cluster name '%s'\n", j.ClusterName)
	fmt.Println(j.Obj)
}
