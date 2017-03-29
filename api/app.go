// mystack api
// https://github.com/topfreegames/mystack
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

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack/mystack-cli/errors"
	"github.com/topfreegames/mystack/mystack-cli/metadata"
	"github.com/topfreegames/mystack/mystack-cli/models"
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
}

//NewApp ctor
func NewApp(host string, port int, debug bool, logger logrus.FieldLogger, config *viper.Viper) (*App, error) {
	controllerProtocol := config.GetString("controller.protocol")
	controllerHost := config.GetString("controller.host")
	controllerPort := config.GetString("controller.port")

	a := &App{
		Address: fmt.Sprintf("%s:%d", host, port),
		Debug:   debug,
		Logger:  logger,
		Login:   models.NewLogin(controllerProtocol, controllerHost, controllerPort),
	}
	err := a.configureApp()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) getRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/healthcheck", Chain(
		&HealthcheckHandler{App: a},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("healthcheck")

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

	err = a.Login.Perform()
	if err != nil {
		return nil, err
	}

	err = a.Server.Serve(listener)
	if err != nil {
		listener.Close()
		return nil, err
	}

	return listener, nil
}
