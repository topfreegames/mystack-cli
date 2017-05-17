// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/topfreegames/mystack-cli/errors"
	"github.com/topfreegames/mystack-cli/mocks"
	"github.com/topfreegames/mystack-cli/models"
)

var _ = Describe("GoogleCallback", func() {
	Describe("SaveAccessToken", func() {
		var (
			expectedState = "state"
			state         = expectedState
			code          = "code"
			env           = "test"
			controllerURL = "http://controller.example.com"
			hosts         = map[string]string{"controller": "controller.mystack.com"}
			fs            = models.NewMockFS(nil)
			url           = fmt.Sprintf("%s/access?code=%s", controllerURL, code)
		)

		It("should create local config file", func() {
			resp, _ := json.Marshal(map[string]interface{}{"token": "i-am-a-token"})
			client := mocks.NewMockClientInterface(mockCtrl)
			client.EXPECT().
				Get(url, hosts["controller"]).
				Return(resp, http.StatusOK, nil)

			err := models.SaveAccessToken(state, code, expectedState, env, controllerURL, hosts, fs, client)
			Expect(err).NotTo(HaveOccurred())

			exists, err := afero.DirExists(fs.AppFS, mystackDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			mystackConfig := filepath.Join(mystackDir, "mystack-test.json")
			exists, err = afero.Exists(fs.AppFS, mystackConfig)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should error if states are not the same", func() {
			state := "unexpeted-state"
			client := mocks.NewMockClientInterface(mockCtrl)

			err := models.SaveAccessToken(state, code, expectedState, env, controllerURL, hosts, fs, client)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("GoogleCallback could not authenticate due to: invalid oauth state, expected 'state', got 'unexpeted-state'"))
		})

		It("should error status is 400", func() {
			client := mocks.NewMockClientInterface(mockCtrl)
			client.EXPECT().
				Get(url, hosts["controller"]).
				Return([]byte("error: bad request"), http.StatusBadRequest, nil)

			err := models.SaveAccessToken(state, code, expectedState, env, controllerURL, hosts, fs, client)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("GoogleCallback could not authenticate due to: error: bad request"))

			oauthError := err.(*errors.OAuthError)
			Expect(oauthError.Status).To(Equal(http.StatusBadRequest))
		})

		It("should error if Get returns error", func() {
			client := mocks.NewMockClientInterface(mockCtrl)
			client.EXPECT().
				Get(url, hosts["controller"]).
				Return(nil, 0, fmt.Errorf("error"))

			err := models.SaveAccessToken(state, code, expectedState, env, controllerURL, hosts, fs, client)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("error"))
		})
	})

	Describe("NewServerControl", func() {
		It("should close listener with channel signal", func() {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "I am a mock")
			}))
			defer mockServer.Close()

			serverControl := models.NewServerControl(mockServer.Listener)
			serverControl.CloseServer <- true

			err := mockServer.Listener.Close()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("use of closed network connection"))
		})
	})
})
