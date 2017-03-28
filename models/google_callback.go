// mystack api
// https://github.com/topfreegames/mystack
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
	"github.com/topfreegames/mystack/mystack-cli/errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

//ServerControl starts a go routine to end the server after callback
type ServerControl struct {
	listener    net.Listener
	CloseServer chan bool
}

func (s *ServerControl) waitClose() error {
	select {
	case <-s.CloseServer:
		err := s.listener.Close()
		return err
	}
}

//NewServerControl is the ServerControl ctor
func NewServerControl(listener net.Listener) *ServerControl {
	serverControl := &ServerControl{
		listener:    listener,
		CloseServer: make(chan bool),
	}

	go serverControl.waitClose()

	return serverControl
}

var (
	tokenDirectory = filepath.Join(os.Getenv("HOME"), ".mystack")
	tokenFile      = "token"
)

func writeFile(str string) error {
	err := os.MkdirAll(tokenDirectory, os.ModePerm)
	if err != nil {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("Couldn't save token due to: '%s'", err))
		return err
	}

	tokenPath := fmt.Sprintf("%s/%s", tokenDirectory, tokenFile)
	bts := []byte(str)
	err = ioutil.WriteFile(tokenPath, bts, 0644)
	if err != nil {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("Couldn't save token due to: '%s'", err))
		return err
	}

	return nil
}

//SaveAccessToken get access token from authorization code and saves locally
func SaveAccessToken(state, code, expectedState string) error {
	if state != expectedState {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("invalid oauth state, expected '%s', got '%s'", expectedState, state))
		return err
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("Code exchange failed with '%s'", err))
		return err
	}
	if !token.Valid() {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("Invalid token received from Authorization Server"))
		return err
	}

	err = writeFile(token.AccessToken)
	return err
}
