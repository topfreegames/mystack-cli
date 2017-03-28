package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Healthcheck", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	It("should receive status code 200 and healthy true", func() {
		request, _ := http.NewRequest("GET", "/healthcheck", nil)
		app.Router.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(recorder.Body.String()).To(Equal("{\"healthy\": true}"))
	})

	It("should receive status code 404 on wrong path", func() {
		request, _ := http.NewRequest("GET", "/healthcheckwrongpath", nil)
		app.Router.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusNotFound))
	})
})
