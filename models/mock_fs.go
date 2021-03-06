// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"os"

	"github.com/spf13/afero"
)

//MockFS implements FileSystem interface
type MockFS struct {
	AppFS afero.Fs
	err   error
}

//NewMockFS constructs a new mock
func NewMockFS(err error) *MockFS {
	return &MockFS{
		AppFS: afero.NewMemMapFs(),
		err:   err,
	}
}

//MkdirAll creates a mock directory
func (m *MockFS) MkdirAll(path string, perm os.FileMode) error {
	if m.err != nil {
		return m.err
	}
	return m.AppFS.MkdirAll(path, perm)
}

//Create creates a mock file
func (m *MockFS) Create(name string) (afero.File, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.AppFS.Create(name)
}

//RemoveAll removes a file
func (m *MockFS) RemoveAll(path string) error {
	if m.err != nil {
		return m.err
	}
	return m.AppFS.RemoveAll(path)
}

//IsNotExist returns true if err if of type FileNotExists
func (m *MockFS) IsNotExist(err error) bool {
	return err != nil
}

//Stat returns the FileInfo describing the the file
func (m *MockFS) Stat(name string) (os.FileInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.AppFS.Stat(name)
}
