// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"github.com/topfreegames/mystack-cli/metadata"
	"net/http"
)

// VersionMiddleware adds the version to the request
type VersionMiddleware struct {
	Next http.Handler
}

//ServeHTTP method
func (m *VersionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Offers-Version", metadata.Version)
	m.Next.ServeHTTP(w, r)
}

//SetNext handler
func (m *VersionMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
