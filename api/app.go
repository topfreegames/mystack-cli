// kubecos api
// https://github.com/topfreegames/kubecos
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
	"github.com/satori/go.uuid"
	"github.com/topfreegames/kubecos/kubecos-cli/errors"
	"github.com/topfreegames/kubecos/kubecos-cli/metadata"
	"github.com/topfreegames/kubecos/kubecos-cli/models"
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
	OAuthState    string
	ServerControl *models.ServerControl
}

//NewApp ctor
func NewApp(host string, port int, debug bool, logger logrus.FieldLogger) (*App, error) {
	a := &App{
		Address:    fmt.Sprintf("%s:%d", host, port),
		Debug:      debug,
		Logger:     logger,
		OAuthState: uuid.NewV4().String(),
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
		"source":    "kubecos-cli",
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

//ListenAndServe requests
func (a *App) ListenAndServe(fn func(...interface{}) error, args ...interface{}) (io.Closer, error) {
	listener, err := net.Listen("tcp", a.Address)
	if err != nil {
		return nil, err
	}

	a.ServerControl = models.NewServerControl(listener)

	err = fn(args...)
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
