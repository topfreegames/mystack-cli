// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//ClientInterface interface
type ClientInterface interface {
	Get(url, host string) ([]byte, int, error)
	GetToStdOut(url, host string) error
	Put(url string, body map[string]interface{}) ([]byte, int, error)
	Delete(url string) ([]byte, int, error)
}
