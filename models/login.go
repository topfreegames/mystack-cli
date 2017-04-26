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
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/satori/go.uuid"
)

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

//Login gets an authorization code from google
type Login struct {
	OAuthState    string
	ControllerURL string
}

//NewLogin is the Login ctor
func NewLogin(controllerURL string) *Login {
	return &Login{
		OAuthState:    randToken(),
		ControllerURL: controllerURL,
	}
}

func randToken() string {
	return uuid.NewV4().String()
}

//Perform makes a request to googleapis
func (l *Login) Perform() (map[string]string, error) {
	path := fmt.Sprintf("%s/login?state=%s", l.ControllerURL, l.OAuthState)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Host = "login"

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d when GET request to controller server", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var bodyObj map[string]string
	json.Unmarshal(body, &bodyObj)
	url := bodyObj["url"]

	err = open(url)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"controller": bodyObj["controllerHost"],
		"logger":     bodyObj["loggerHost"],
	}, nil
}
