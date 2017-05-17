// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/spf13/afero"
	"os"
)

//RealFS implements FileSystem interface
type RealFS struct{}

//NewRealFS constructs a new mock
func NewRealFS() *RealFS {
	return &RealFS{}
}

//MkdirAll creates a mock directory
func (m *RealFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

//Create creates a mock file
func (m *RealFS) Create(name string) (afero.File, error) {
	return os.Create(name)
}

//RemoveAll removes a file
func (m *RealFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

//IsNotExist returns true if err if of type FileNotExists
func (m *RealFS) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

//Stat returns the FileInfo describing the the file
func (m *RealFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
