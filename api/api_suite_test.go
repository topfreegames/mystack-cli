package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/mystack/mystack-cli/api"
	"github.com/topfreegames/mystack/mystack-cli/models"
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

	var err error
	app, err = api.NewApp("0.0.0.0", 8989, false, l, models.NewLogin())
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
})
