// mystack api
// https://github.com/topfreegames/mystack
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"
)

//HealthcheckHandler handler
type HealthcheckHandler struct {
	App *App
}

//ServeHTTP method
func (h *HealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := loggerFromContext(r.Context())

	l.Debug("Performing healthcheck...")

	Write(w, http.StatusOK, `{"healthy": true}`)
	l.Debug("Healthcheck done.")
}
