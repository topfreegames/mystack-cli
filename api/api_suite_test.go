// mystack-cli api
// https://github.com/topfreegames/mystack-cli //
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/mystack-cli/api"
	"testing"
)

var app *api.App

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var _ = BeforeSuite(func() {
	l := logrus.New()
	l.Level = logrus.FatalLevel

	_, err := api.NewApp("0.0.0.0", 57459, false, l, "production", "http://localhost:8080", "mystack.controller.com")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
})
