package models_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	. "github.com/topfreegames/mystack-cli/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Proxy", func() {
	getPort := func() (int, error) {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return 0, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return 0, err
		}
		defer l.Close()
		return l.Addr().(*net.TCPAddr).Port, nil
	}

	url := func(port int) string {
		return fmt.Sprintf("localhost:%d", port)
	}

	Describe("Start", func() {
		It("should start proxy", func() {
			from, err := getPort()
			Expect(err).NotTo(HaveOccurred())
			to, err := getPort()
			Expect(err).NotTo(HaveOccurred())

			afterMessage := "message"
			proxy := NewProxy(url(from), url(to), map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			})

			close := make(chan bool)
			ln, err := net.Listen("tcp", url(to))
			Expect(err).NotTo(HaveOccurred())
			go func() {
				conn, err := ln.Accept()
				Expect(err).NotTo(HaveOccurred())
				defer func() {
					select {
					case <-close:
						conn.Close()
					}
				}()

				buf, err := bufio.NewReader(conn).ReadBytes('\n')
				Expect(err).NotTo(HaveOccurred())
				conn.Write([]byte("ok"))

				message := make(map[string]interface{})
				err = json.Unmarshal(buf, &message)
				Expect(err).NotTo(HaveOccurred())
				Expect(message).To(HaveKeyWithValue("key1", "value1"))
				Expect(message).To(HaveKeyWithValue("key2", "value2"))

				buf = make([]byte, len(afterMessage))
				_, err = conn.Read(buf)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(buf)).To(Equal(afterMessage))
				conn.Write([]byte("ok"))
			}()

			go proxy.Start()
			defer proxy.Stop()
			conn, err := net.Dial("tcp", url(from))
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(conn, afterMessage)
			buf := make([]byte, 2)
			close <- true
			conn.Read(buf)
			Expect(string(buf)).To(Equal("ok"))
		})
	})
})
