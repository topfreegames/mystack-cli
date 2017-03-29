// mystack-cli api
// https://github.com/topfreegames/mystack/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"github.com/topfreegames/mystack/mystack-cli/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {
	Describe("Login Function", func() {
		It("should not return error on calling Login with string", func() {
			login := models.NewLogin("http://localhost:8888")
			err := login.Perform()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
