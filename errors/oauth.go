// kubecos api
// https://github.com/topfreegames/kubecos
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package errors

import (
	"encoding/json"
	"fmt"
)

//OAuthError occurs during authentication phase
type OAuthError struct {
	Model   string
	Message string
}

//NewOAuthError ctor
func NewOAuthError(model, message string) *OAuthError {
	return &OAuthError{
		Model:   model,
		Message: message,
	}
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("%s could not authenticate due to: %s", e.Model, e.Message)
}

//Serialize returns the error serialized
func (e *OAuthError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-002",
		"error":       fmt.Sprintf("Conflicted%sError", e.Model),
		"description": e.Error(),
	})

	return g
}
