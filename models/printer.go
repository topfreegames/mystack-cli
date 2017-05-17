// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//Printer is the printer interface
//Helps unifying printers that print in different ways
type Printer interface {
	Print()
}
