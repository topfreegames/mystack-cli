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
	"github.com/topfreegames/mystack-cli/errors"
	"io/ioutil"
	"net"
	"net/http"
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
func SaveAccessToken(state, code, expectedState, env, controllerURL, controllerHost string) error {
	if state != expectedState {
		err := errors.NewOAuthError("GoogleCallback", fmt.Sprintf("invalid oauth state, expected '%s', got '%s'", expectedState, state))
		return err
	}

	url := fmt.Sprintf("%s/access?code=%s", controllerURL, code)
	req, err := http.NewRequest("GET", url, nil)
	req.Host = controllerHost

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var bodyObj map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Status: %d\nError: %s", resp.StatusCode, string(body))
		return err
	}
	json.Unmarshal(body, &bodyObj)
	token := bodyObj["token"].(string)

	c := NewConfig(env, token, controllerURL, controllerHost)
	err = c.Write()
	return err
}
