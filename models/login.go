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
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
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
	OAuthState string
	ServerURL  string
}

//NewLogin is the Login ctor
func NewLogin(controllerURL string) *Login {
	return &Login{
		OAuthState: randToken(),
		ServerURL:  controllerURL,
	}
}

func randToken() string {
	return uuid.NewV4().String()
}

//Perform makes a request to googleapis
func (l *Login) Perform() error {
	basePath := l.ServerURL
	resp, err := http.Get(fmt.Sprintf("%s/login?state=%s", basePath, l.OAuthState))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Status code %d when GET request to controller server", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var bodyObj map[string]interface{}
	json.Unmarshal(body, &bodyObj)
	url := bodyObj["url"].(string)

	err = open(url)

	return err
}
