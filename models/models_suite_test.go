// mystack-cli api
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/user"
	"path/filepath"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var mystackDir, responseHost, responseAuth string
var mockCtrl *gomock.Controller
var hosts map[string]string
var mockServer *httptest.Server

const msg, controllerHost string = "I am a mock", "controller.example.com"

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	usr, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	mystackDir = filepath.Join(usr.HomeDir, ".mystack")
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseHost = r.Host
		responseAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, msg)
	}))
})

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
})

var _ = AfterSuite(func() {
	mockServer.Close()
})
