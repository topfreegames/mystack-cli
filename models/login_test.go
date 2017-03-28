// kubecos api
// https://github.com/topfreegames/kubecos
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/kubecos/kubecos-cli/models"
)

var _ = Describe("Login Model", func() {
	Describe("Login", func() {
		It("should not return error on calling Login with string", func() {
			state := "random-state"
			err := models.Login(state)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error on calling Login with integer", func() {
			state := 123
			err := models.Login(state)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
