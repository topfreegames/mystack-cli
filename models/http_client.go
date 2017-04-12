// mystack
// https://github.com/topfreegames/mystack-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	h.client = &http.Client{}
	return h
}

// GET does a get request
func (c *MyStackHTTPClient) GET(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	auth := fmt.Sprintf("Bearer %s", c.config.Token)
	req.Header.Add("Authorization", auth)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
