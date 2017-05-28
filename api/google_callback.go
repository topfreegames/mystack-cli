// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"io"
	"net/http"

	"github.com/topfreegames/mystack-cli/errors"
	"github.com/topfreegames/mystack-cli/models"
)

const Index = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Mystack</title>
</head>
<body>
  <h1>Thanks for logging in</h1>
  You can go back to your terminal
</body>
</html>
`

const UnauthorizedIndex = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Mystack</title>
</head>
<body>
  <h1>Unauthorized</h1>
  Your email is not authorized to use Mystack
</body>
</html>
`

//OAuthCallbackHandler handles the callback after user approves/deny auth
type OAuthCallbackHandler struct {
	app    *App
	fs     models.FileSystem
	client models.ClientInterface
}

func NewOAuthCallbackHandler(
	app *App,
	fs models.FileSystem,
	client models.ClientInterface,
) *OAuthCallbackHandler {
	return &OAuthCallbackHandler{
		app:    app,
		fs:     fs,
		client: client,
	}
}

//ServeHTTP method
func (o *OAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	l := loggerFromContext(r.Context())
	l.Debugf("Returned state %s and code %s", state, code)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := models.SaveAccessToken(
		state, code, o.app.Login.OAuthState, o.app.env, o.app.Login.ControllerURL,
		o.app.Hosts,
		o.fs,
		o.client,
	)
	if err != nil {
		if _, ok := err.(*errors.OAuthError); ok {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, UnauthorizedIndex)
			o.app.ServerControl.CloseServer <- true
			return
		}

		l.Error(err)
		o.app.HandleError(w, http.StatusInternalServerError, "unexpected authentication error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, Index)

	o.app.ServerControl.CloseServer <- true
}
