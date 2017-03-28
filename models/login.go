// mystack api
// https://github.com/topfreegames/mystack
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"os/exec"
	"runtime"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("googlekey"),
		ClientSecret: os.Getenv("googlesecret"),
		RedirectURL:  "http://localhost:8989/google-callback",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
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
}

//NewLogin is the Login ctor
func NewLogin() *Login {
	return &Login{
		OAuthState: uuid.NewV4().String(),
	}
}

//Perform makes a request to googleapis
func (l *Login) Perform() error {
	url := googleOauthConfig.AuthCodeURL(l.OAuthState)
	err := open(url)

	return err
}
