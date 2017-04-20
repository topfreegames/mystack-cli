// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"net/http"

	"github.com/topfreegames/mystack-cli/errors"
	"github.com/topfreegames/mystack-cli/models"
)

const index = `
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

//OAuthCallbackHandler handles the callback after user approves/deny auth
type OAuthCallbackHandler struct {
	App *App
}

//ServeHTTP method
func (o *OAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	l := loggerFromContext(r.Context())

	err := models.SaveAccessToken(
		state,
		code,
		o.App.Login.OAuthState,
		o.App.env,
		o.App.Login.ServerURL,
		o.App.Login.ServerHost,
	)
	if err != nil {
		if err, ok := err.(*errors.OAuthError); ok {
			l.Error(err.Serialize())
		}

		l.Error(err)
		return
	}

	o.App.ServerControl.CloseServer <- true
	fmt.Fprintf(w, index)
}
