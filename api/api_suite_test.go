// mystack-cli api
// https://github.com/topfreegames/mystack-cli //
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

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

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/mystack-cli/api"
	"github.com/topfreegames/mystack-cli/models"
)

var app *api.App
var mystackDir string
var mockCtrl *gomock.Controller
var hosts map[string]string
var mockServer *httptest.Server

const controllerHost string = "http://controller.example.com"

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var _ = BeforeSuite(func() {
	l := logrus.New()
	l.Level = logrus.FatalLevel

	usr, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	mystackDir = filepath.Join(usr.HomeDir, ".mystack")

	app, err = api.NewApp("0.0.0.0", 57459, false, l, "test", controllerHost)
	Expect(err).NotTo(HaveOccurred())

	hosts = map[string]string{"controller": controllerHost}
	app.Hosts = hosts
})

var _ = AfterSuite(func() {
})

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "I am a mock")
	}))
	app.ServerControl = models.NewServerControl(mockServer.Listener)
})

var _ = AfterEach(func() {
	mockServer.Close()
})
