package models_test

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/topfreegames/mystack-cli/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpClient", func() {
	hosts = map[string]string{"controller": controllerHost}
	config := NewConfig("test", "token", controllerHost, hosts)
	client := NewMyStackHTTPClient(config)

	Describe("Get", func() {
		It("should get with correct headers", func() {
			url := fmt.Sprintf("http://%s", mockServer.Listener.Addr().String())
			resp, status, err := client.Get(url, controllerHost)
			Expect(err).NotTo(HaveOccurred())

			respStr := strings.TrimSpace(string(resp))
			Expect(respStr).To(Equal(msg))
			Expect(status).To(Equal(http.StatusOK))
			Expect(responseAuth).To(Equal("Bearer token"))
			Expect(responseHost).To(Equal(controllerHost))
		})

		It("should return error with invalid url", func() {
			_, _, err := client.Get("qwerty://iaminvalid", controllerHost)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Put", func() {
		It("should put with correct headers", func() {
			url := fmt.Sprintf("http://%s", mockServer.Listener.Addr().String())
			resp, status, err := client.Put(url, map[string]interface{}{"key": "value"})
			Expect(err).NotTo(HaveOccurred())

			respStr := strings.TrimSpace(string(resp))
			Expect(respStr).To(Equal(msg))
			Expect(status).To(Equal(http.StatusOK))
			Expect(responseAuth).To(Equal("Bearer token"))
			Expect(responseHost).To(Equal(controllerHost))
		})

		It("should put with correct headers with nil body", func() {
			url := fmt.Sprintf("http://%s", mockServer.Listener.Addr().String())
			resp, status, err := client.Put(url, nil)
			Expect(err).NotTo(HaveOccurred())

			respStr := strings.TrimSpace(string(resp))
			Expect(respStr).To(Equal(msg))
			Expect(status).To(Equal(http.StatusOK))
			Expect(responseAuth).To(Equal("Bearer token"))
			Expect(responseHost).To(Equal(controllerHost))
		})

		It("should return error with invalid url", func() {
			_, _, err := client.Put("qwerty://iaminvalid", nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Delete", func() {
		It("should delete with correct headers", func() {
			url := fmt.Sprintf("http://%s", mockServer.Listener.Addr().String())
			resp, status, err := client.Delete(url)
			Expect(err).NotTo(HaveOccurred())

			respStr := strings.TrimSpace(string(resp))
			Expect(respStr).To(Equal(msg))
			Expect(status).To(Equal(http.StatusOK))
			Expect(responseAuth).To(Equal("Bearer token"))
			Expect(responseHost).To(Equal(controllerHost))
		})

		It("should return error with invalid url", func() {
			_, _, err := client.Delete("qwerty://iaminvalid")
			Expect(err).To(HaveOccurred())
		})
	})
})
