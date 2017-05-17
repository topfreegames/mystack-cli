// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/topfreegames/mystack-cli/errors"
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

//SaveAccessToken get access token from authorization code and saves locally
func SaveAccessToken(
	state, code, expectedState, env, controllerURL string,
	hosts map[string]string,
	fs FileSystem,
	client ClientInterface,
) error {
	if state != expectedState {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("invalid oauth state, expected '%s', got '%s'", expectedState, state), http.StatusBadRequest)
		return err
	}

	url := fmt.Sprintf("%s/access?code=%s", controllerURL, code)
	resp, status, err := client.Get(url, hosts["controller"])
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.NewOAuthError("GoogleCallback", string(resp), status)
	}

	var bodyObj map[string]interface{}
	json.Unmarshal(resp, &bodyObj)
	token := bodyObj["token"].(string)

	c := NewConfig(env, token, controllerURL, hosts)
	err = c.Write(fs)
	return err
}
