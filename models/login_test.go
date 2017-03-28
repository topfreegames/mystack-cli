package models_test

import (
	"github.com/topfreegames/kubecos/kubecos-cli/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {
	Describe("Login Function", func() {
		It("should not return error on calling Login with string", func() {
			login := models.NewLogin()
			err := login.Perform()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
