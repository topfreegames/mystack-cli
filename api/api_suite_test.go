package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/mystack/mystack-cli/api"
	mTesting "github.com/topfreegames/mystack/mystack-cli/testing"
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

	config, err := mTesting.GetDefaultConfig()
	Expect(err).NotTo(HaveOccurred())

	app, err = api.NewApp("0.0.0.0", 8989, false, l, config)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
})
