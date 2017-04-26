// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/topfreegames/mystack-cli/errors"
	"github.com/topfreegames/mystack-cli/metadata"
	"github.com/topfreegames/mystack-cli/models"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

//App is our API application
type App struct {
	Address       string
	DB            runner.Connection
	Debug         bool
	Logger        logrus.FieldLogger
	Router        *mux.Router
	Server        *http.Server
	ServerControl *models.ServerControl
	Login         *models.Login
	env           string
	controllerURL string
	Hosts         map[string]string
}

//NewApp ctor
func NewApp(host string, port int, debug bool, logger logrus.FieldLogger, env, controllerURL string) (*App, error) {
	a := &App{
		Address: fmt.Sprintf("%s:%d", host, port),
		Debug:   debug,
		Logger:  logger,
		Login:   models.NewLogin(controllerURL),
		env:     env,
	}
	err := a.configureApp()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) getRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/google-callback", Chain(
		&OAuthCallbackHandler{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("oauth2")

	return r
}

func (a *App) configureApp() error {
	a.configureLogger()
	a.configureServer()
	return nil
}

func (a *App) configureLogger() {
	a.Logger = a.Logger.WithFields(logrus.Fields{
		"source":    "mystack-cli",
		"operation": "initializeApp",
		"version":   metadata.Version,
	})
}

func (a *App) configureServer() {
	a.Router = a.getRouter()
	a.Server = &http.Server{Addr: a.Address, Handler: a.Router}
}

//HandleError writes an error response with message and status
func (a *App) HandleError(w http.ResponseWriter, status int, msg string, err interface{}) {
	w.WriteHeader(status)
	var sErr errors.SerializableError
	val, ok := err.(errors.SerializableError)
	if ok {
		sErr = val
	} else {
		sErr = errors.NewGenericError(msg, err.(error))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(sErr.Serialize())
}

//ListenAndLoginAndServe logins and starts local server to get access token from google
func (a *App) ListenAndLoginAndServe() (io.Closer, error) {
	listener, err := net.Listen("tcp", a.Address)
	if err != nil {
		return nil, err
	}

	a.ServerControl = models.NewServerControl(listener)

	hosts, err := a.Login.Perform()
	if err != nil {
		return nil, err
	}

	a.Hosts = hosts

	err = a.Server.Serve(listener)
	//TODO: do a better check, in case a real "use of closed network connection" happens
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		listener.Close()
		return nil, err
	}

	return listener, nil
}
