// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	. "github.com/topfreegames/mystack-cli/api"
	"github.com/topfreegames/mystack-cli/metadata"
	"github.com/topfreegames/mystack-cli/mocks"
	"github.com/topfreegames/mystack-cli/models"
)

var _ = Describe("GoogleCallback", func() {
	Describe("GET /google-callback", func() {
		var response *httptest.ResponseRecorder

		BeforeEach(func() {
			response = httptest.NewRecorder()
		})

		var (
			code          = "code"
			state         = "state"
			fs            = models.NewMockFS(nil)
			controllerURL = fmt.Sprintf("%s/access?code=%s", controllerHost, code)
		)

		It("should create local config file", func() {
			resp, _ := json.Marshal(map[string]interface{}{"token": app.Login.OAuthState})
			client := mocks.NewMockClientInterface(mockCtrl)
			client.EXPECT().
				Get(controllerURL, hosts["controller"]).
				Return(resp, http.StatusOK, nil).AnyTimes()

			o := NewOAuthCallbackHandler(app, fs, client)

			request, _ := http.NewRequest("GET", "http://localhost:53789", nil)
			q := request.URL.Query()
			q.Add("state", app.Login.OAuthState)
			q.Add("code", code)
			request.URL.RawQuery = q.Encode()

			logMid := &LoggingMiddleware{App: app, Next: o}
			versionMid := &VersionMiddleware{Next: logMid}
			versionMid.ServeHTTP(response, request)

			Expect(response.Header().Get("Content-Type")).To(Equal("text/html; charset=utf-8"))
			Expect(response.Header().Get("X-Mystack-Version")).To(Equal(metadata.Version))
			Expect(response.Body.String()).To(Equal(Index))

			exists, err := afero.DirExists(fs.AppFS, mystackDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			mystackConfig := filepath.Join(mystackDir, "mystack-test.json")
			exists, err = afero.Exists(fs.AppFS, mystackConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should error if states are not the same", func() {
			client := mocks.NewMockClientInterface(mockCtrl)

			o := NewOAuthCallbackHandler(app, fs, client)

			request, _ := http.NewRequest("GET", "http://localhost:53789", nil)
			q := request.URL.Query()
			q.Add("state", state)
			q.Add("code", code)
			request.URL.RawQuery = q.Encode()

			logMid := &LoggingMiddleware{App: app, Next: o}
			logMid.ServeHTTP(response, request)

			Expect(response.Header().Get("Content-Type")).To(Equal("text/html; charset=utf-8"))
			Expect(response.Body.String()).To(Equal(UnauthorizedIndex))
		})

		It("should error status is 400", func() {
			client := mocks.NewMockClientInterface(mockCtrl)
			client.EXPECT().
				Get(controllerURL, hosts["controller"]).
				Return([]byte("error: bad request"), http.StatusBadRequest, nil)

			o := NewOAuthCallbackHandler(app, fs, client)

			request, _ := http.NewRequest("GET", "http://localhost:53789", nil)
			q := request.URL.Query()
			q.Add("state", app.Login.OAuthState)
			q.Add("code", code)
			request.URL.RawQuery = q.Encode()

			logMid := &LoggingMiddleware{App: app, Next: o}
			logMid.ServeHTTP(response, request)

			Expect(response.Header().Get("Content-Type")).To(Equal("text/html; charset=utf-8"))
			Expect(response.Body.String()).To(Equal(UnauthorizedIndex))
		})
	})
})
