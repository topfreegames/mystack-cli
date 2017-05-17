// mystack-cli
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// MyStackHTTPClient struct
type MyStackHTTPClient struct {
	config *Config
	client *http.Client
}

// NewMyStackHTTPClient ctor
func NewMyStackHTTPClient(config *Config) *MyStackHTTPClient {
	h := &MyStackHTTPClient{
		config: config,
	}

	h.client = &http.Client{
		Timeout: 20 * time.Minute,
	}

	return h
}

// Get does a get request
func (c *MyStackHTTPClient) Get(url, host string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	c.addAuthHeader(req, host)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return body, res.StatusCode, nil
}

// GetToStdOut streams the content received from server to stdout
func (c *MyStackHTTPClient) GetToStdOut(url, host string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	c.addAuthHeader(req, host)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		fmt.Printf(string(line))
	}
}

// Put does a put request
func (c *MyStackHTTPClient) Put(url string, body map[string]interface{}) ([]byte, int, error) {
	ioBody, err := ioReader(body)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("PUT", url, ioBody)
	if err != nil {
		return nil, 0, err
	}
	req.Close = true

	c.addAuthHeader(req, c.config.ControllerHost)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("Error creating cluster")
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return responseBody, res.StatusCode, nil
}

func (c *MyStackHTTPClient) addAuthHeader(req *http.Request, host string) {
	if c.config != nil {
		auth := fmt.Sprintf("Bearer %s", c.config.Token)
		req.Header.Add("Authorization", auth)
	}

	if host != "" {
		req.Host = host
	}
}

func ioReader(body map[string]interface{}) (*bytes.Reader, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}

// Delete does a put request
func (c *MyStackHTTPClient) Delete(url string) ([]byte, int, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, 0, err
	}

	c.addAuthHeader(req, c.config.ControllerHost)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return responseBody, res.StatusCode, nil
}
