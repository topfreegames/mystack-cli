package api

import (
	"github.com/topfreegames/kubecos/kubecos-cli/models"
	"net/http"
)

//OAuthCallbackHandler handles the callback after user approves/deny auth
type OAuthCallbackHandler struct {
	App *App
}

//ServeHTTP method
func (o *OAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	_ = models.SaveAccessToken(state, code, o.App.OAuthState)
	//TODO: return error code

	o.App.Listener.Close()
}
